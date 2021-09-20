package lib

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/e421083458/gorm"
	_ "github.com/e421083458/gorm/dialects/mysql"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

func InitDBPool(path string) error {
	DbConfMap := &MysqlMapConf{}
	err := ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}

	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for confName, DbConf := range DbConfMap.List {
		dbpool, err := sql.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbpool.SetMaxOpenConns(DbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		err = dbpool.Ping()
		if err != nil {
			return err
		}
		dbgorm, err := gorm.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbgorm.SingularTable(true)
		err = dbgorm.DB().Ping()
		if err != nil {
			return err
		}
		dbgorm.LogMode(true)
		dbgorm.LogCtx(true)
		dbgorm.SetLogger(&MysqlGormLogger{})
		dbgorm.DB().SetMaxIdleConns(DbConf.MaxIdleConn)
		dbgorm.DB().SetMaxOpenConns(DbConf.MaxOpenConn)
		dbgorm.DB().SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		DBMapPool[confName] = dbpool
		GORMMapPool[confName] = dbgorm
	}
	return nil
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func CloseDB() error {
	for _, dbpool := range DBMapPool {
		dbpool.Close()
	}
	for _, dbpool := range GORMMapPool {
		dbpool.Close()
	}
	return nil
}

type MysqlGormLogger struct {
	gorm.Logger
	Trace *trace.Trace
}

func (logger *MysqlGormLogger) CtxPrint(s *gorm.DB, values ...interface{}) {
	ctx, ok := s.GetCtx()
	var traceContext *trace.Trace
	if ok {
		traceContext = ctx.(*trace.Trace)
	}
	if !ok {
		traceContext = trace.New(nil)
	}
	message := logger.LogFormatter(values...)
	if message["level"] == "sql" {
		Log.TagInfo(traceContext, trace.DLTagMysqlSuccess, message)
	} else {
		Log.TagWarn(traceContext, trace.DLTagMysqlFailed, message)
	}
}

func (logger *MysqlGormLogger) LogFormatter(values ...interface{}) (messages map[string]interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = logger.NowFunc().Format("2006-01-02 15:04:05")
			source          = fmt.Sprintf("%v", values[1])
		)

		messages = map[string]interface{}{"level": level, "source": source, "current_time": currentTime}

		if level == "sql" {
			messages["proc_time"] = fmt.Sprintf("%fs", values[2].(time.Duration).Seconds())
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); logger.isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if regexp.MustCompile(`\$\d+`).MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range regexp.MustCompile(`\?`).Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}
			messages["sql"] = sql
			if len(values) > 5 {
				messages["affected_row"] = strconv.FormatInt(values[5].(int64), 10)
			}
		} else {
			messages["ext"] = values
		}
	}
	return
}

func (logger *MysqlGormLogger) NowFunc() time.Time {
	return time.Now()
}

func (logger *MysqlGormLogger) isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
