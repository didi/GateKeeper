package lib

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/didi/gatekeeper/golang_common/trace"
	"github.com/didi/gatekeeper/golang_common/zerolog/ddlog"
	"github.com/didi/gatekeeper/golang_common/zerolog/log"
	"github.com/e421083458/gorm"
	"github.com/spf13/viper"
	"io/ioutil"
	//"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	ConfBase     *BaseConf
	DBMapPool    map[string]*sql.DB
	GORMMapPool  map[string]*gorm.DB
	ConfRedis    *RedisConf
	ConfRedisMap *RedisMapConf
	ViperConfMap map[string]*viper.Viper
	TimeLocation *time.Location
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
	LocalIP      = net.ParseIP("127.0.0.1")
	ConfEnvPath  string
	ConfEnv      string
	ZLog         *ddlog.DiLogHandle
)

type BaseConf struct {
	TimeLocation string        `toml:"time_location"`
	Log          ZeroLogConfig `toml:"log"`
	Base         struct {
		DebugMode    string `toml:"debug_mode"`
		TimeLocation string `toml:"time_location"`
		SerName      string `toml:"ser_name"`
	} `toml:"base"`
}

type LogConfFileWriter struct {
	On              bool   `toml:"on"`
	LogPath         string `toml:"log_path"`
	RotateLogPath   string `toml:"rotate_log_path"`
	WfLogPath       string `toml:"wf_log_path"`
	RotateWfLogPath string `toml:"rotate_wf_log_path"`
}

type LogConfConsoleWriter struct {
	On    bool `toml:"on"`
	Color bool `toml:"color"`
}

type LogConfig struct {
	Level string               `toml:"log_level"`
	FW    LogConfFileWriter    `toml:"file_writer"`
	CW    LogConfConsoleWriter `toml:"console_writer"`
}

type ZeroLogConfig struct {
	On          bool   `toml:"on"`
	Level       string `toml:"level"`
	FilePrefix  string `toml:"file_prefix"`
	FileDir     string `toml:"file_dir"`
	AutoClear   bool   `toml:"auto_clear"`
	ClearHours  int32  `toml:"clear_hours"`
	ClearStep   int32  `toml:"clear_step"`
	Separate    bool   `toml:"separate"`
	DisableLink bool   `toml:"disable_link"`
}

type MysqlMapConf struct {
	List map[string]*MySQLConf `toml:"list"`
}

type MySQLConf struct {
	DriverName      string `toml:"driver_name"`
	DataSourceName  string `toml:"data_source_name"`
	MaxOpenConn     int    `toml:"max_open_conn"`
	MaxIdleConn     int    `toml:"max_idle_conn"`
	MaxConnLifeTime int    `toml:"max_conn_life_time"`
}

type RedisMapConf struct {
	List map[string]*RedisConf `toml:"list"`
}

type RedisConf struct {
	ProxyList    []string `toml:"proxy_list"`
	Password     string   `toml:"password"`
	Db           int      `toml:"db"`
	ConnTimeout  int      `toml:"conn_timeout"`
	ReadTimeout  int      `toml:"read_timeout"`
	WriteTimeout int      `toml:"write_timeout"`
}

func ParseConfPath(config string) error {
	if strings.LastIndex(config, "/") != (len(config) - 1) {
		config = config + "/"
	}
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]
	return nil
}

func GetConfEnv() string {
	return ConfEnv
}

func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}

func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}

func ParseLocalConfig(fileName string, st interface{}) error {
	path := GetConfFilePath(fileName)
	err := ParseConfig(path, st)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open config %v fail, %v", path, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read config fail, %v", err)
	}
	if err := toml.Unmarshal(data, conf); err != nil {
		return fmt.Errorf("Parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}

func InitConf(configPath string) error {
	if configPath == "" {
		fmt.Println("input config file like ./conf/dev/")
		os.Exit(1)
	}
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIP = ips[0]
	}
	if err := ParseConfPath(configPath); err != nil {
		return err
	}
	if err := InitViperConf(); err != nil {
		return err
	}
	if err := InitBaseConf(GetConfPath("base")); err != nil {
		log.Error().Msg(Purple(err.Error()))
	}
	if err := InitRedisConf(GetConfPath("redis_map")); err != nil {
		log.Error().Msg(Purple(err.Error()))
	}
	if err := InitDBPool(GetConfPath("mysql_map")); err != nil {
		log.Error().Msg(Purple(err.Error()))
		return err
	}
	if location, err := time.LoadLocation(ConfBase.TimeLocation); err != nil {
		log.Error().Msg(Purple(err.Error()))
		return err
	} else {
		TimeLocation = location
	}
	log.Info().Msg(Purple("success initialized application configuration"))
	return nil
}

func DestroyConf() {
	CloseDB()
	log.Info().Msg(Purple("success destroy config."))
}

func InitBaseConf(path string) error {
	ConfBase = &BaseConf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
		return err
	}
	if ConfBase.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.TimeLocation = "Asia/Chongqing"
		}
	}
	if ConfBase.Log.Level == "" {
		ConfBase.Log.Level = "info"
	}
	config := ddlog.FileConfig{}
	config.Level = ConfBase.Log.Level
	config.FilePrefix = ConfBase.Log.FilePrefix
	config.FileDir = ConfBase.Log.FileDir
	config.AutoClear = ConfBase.Log.AutoClear
	config.ClearHours = ConfBase.Log.ClearHours
	config.ClearStep = ConfBase.Log.ClearStep
	config.Separate = ConfBase.Log.Separate
	config.DisableLink = ConfBase.Log.DisableLink
	tmpLog, err := ddlog.NewLoggerWithCfg(&config)
	if err != nil {
		return err
	}
	ZLog = tmpLog
	ZLog.CtxFormatFunc = trace.FormatCtx
	return nil
}

func InitRedisConf(path string) error {
	ConfRedis := &RedisMapConf{}
	err := ParseConfig(path, ConfRedis)
	if err != nil {
		return err
	}
	ConfRedisMap = ConfRedis
	return nil
}

func InitViperConf() error {
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}

func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
	return confString
}

func GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetConf(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.Get(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:len(keys)], "."))
	return conf
}

func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.IsSet(strings.Join(keys[1:len(keys)], "."))
	return conf
}
