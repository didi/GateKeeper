package check

import (
	"database/sql"
	"fmt"
	"gatekeeper/install/template"
	"gatekeeper/install/tool"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Mysql struct{
	Host 	 string
	Port 	 string
	User     string
	Pwd	 	 string
	Database string
}

var (
	DbPool  *sql.Tx
	err     error
	MysqlClient Mysql
)


func (m Mysql) Init() error{

	// connect mysql
	mysqlLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/", m.User, m.Pwd, m.Host, m.Port)

	db, _ := sql.Open("mysql", mysqlLink)
	if err := db.Ping(); err != nil {
		tool.LogWarning.Println(err)
		return InitDb()
		//return errors.New("connect mysql error")
	}
	// check connect
	db.SetConnMaxLifetime(time.Second * 30)
	DbPool, err = db.Begin()
	if err != nil {
		tool.LogError.Println(err.Error())
		return errors.New("db error")
	}

	tool.LogInfo.Println("connect mysql success")

	tool.LogInfo.Println("init mysql db start")
	// check database
	err = checkDb(m.Database)
	if err != nil{
		tool.LogInfo.Println("init mysql db end")
		DbPool.Rollback()
		return err
	}
	tool.LogInfo.Println("init mysql db end")


	// check table
	err = template.InitSql()
	if err != nil{
		tool.LogInfo.Println("init mysql table end")
		DbPool.Rollback()
		return err
	}
	tool.LogInfo.Println("init mysql table start")
	err = checkTable(m.Database)
	if err != nil{
		tool.LogInfo.Println("init mysql table end")
		DbPool.Rollback()
		return err
	}
	tool.LogInfo.Println("init mysql table end")

	defer DbPool.Commit()
	return nil
}


func InitDb() error{
	host, err := tool.Input("please enter mysql host (default:127.0.0.1):", "127.0.0.1")
	//mysqlLinkInfo, err := inputHost(mysqlLinkInfo)
	if err != nil{
		return err
	}

	port, err := tool.Input("please enter mysql port (default:3306):", "3306")
	//port, err := inputPort(mysqlLinkInfo)
	if err != nil{
		return err
	}

	user, err := tool.Input("please enter mysql user (default:root):", "root")
	if err != nil{
		return err
	}


	pwd, err := tool.Input("please enter mysql pwd (default:root):", "root")
	if err != nil{
		return err
	}

	database, err := tool.Input("please enter database (default:gatekeeper):", "gatekeeper")
	if err != nil{
		return err
	}
	mysql := Mysql{
		Host: host,
		Port: port,
		User: user,
		Pwd: pwd,
		Database: database,
	}
	MysqlClient = mysql
	tool.LogInfo.Println(fmt.Sprintf("mysql connect info host:[%s] port:[%s] user:[%s] pwd:[%s] database[%s]", host, port, user, pwd, database))

	err = mysql.Init();if err !=nil{
		return err
	}
	return nil
}


func checkDb(database string) error {
	dbSql := fmt.Sprintf("USE %s", database)
	tool.LogInfo.Println(dbSql)

	_, err := DbPool.Exec(dbSql)
	if err != nil{
		tool.LogWarning.Println(err.Error())
		// database not exist
		if strings.Contains(err.Error(), "Unknown database") {
			boolCreateDb, err := tool.Confirm("create DB ["+database+"]", 3)
			if err != nil {
				return err
			}

			// create database
			if boolCreateDb {
				return createDb(database)
			} else {
				return errors.New("no database selected")
			}
		}
	}

	_, err = DbPool.Exec("USE " + database); if err != nil{
		return err
	}

	return nil
}


func createDb(database string) error {
	createDbSql := fmt.Sprintf("CREATE DATABASE %s", database)
	tool.LogInfo.Println(createDbSql)
	_, err := DbPool.Exec(createDbSql)
	if err != nil{
		return err
	}
	tool.LogInfo.Println("create database [" + database + "] success")
	return nil
}


func checkTable(database string) error {
	tables := template.Tables
	for table, createTableSql := range tables{
		// check table exist
		//checkTableSql := fmt.Sprintf(
		//	"SELECT COLUMN_NAME fName,column_comment fDesc,DATA_TYPE dataType, " +
		//		"IS_NULLABLE isNull,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) sLength " +
		//		"FROM information_schema.columns " +
		//		"WHERE " +
		//		"table_schema = '%s' " +
		//		"AND table_name = '%s'",
		//	database,
		//	table)
		checkTableSql := fmt.Sprintf("SHOW CREATE TABLE %s.%s", database, table)
		tool.LogInfo.Println("check table [" + table + "]")
		tool.LogInfo.Println(checkTableSql)
		rows, err := DbPool.Query(checkTableSql)
		rows.Close()
		if err != nil{
			// table not exist
			if strings.Contains(err.Error(), "doesn't exist") {
				// create table
				tool.LogInfo.Println("create table [" + table + "]")
				err := createTable(database, createTableSql); if err != nil{
					return err
				}
				continue
			}
			// other err
			return err
		} else{
			// create table
			//type Field struct {
			//	fieldName string
			//	fieldDesc string
			//	dataType  string
			//	isNull    string
			//	length    int
			//}
			//for rows.Next() {
			//	var f Field
			//	err = rows.Scan(&f.fieldName, &f.fieldDesc, &f.dataType, &f.isNull, &f.length)
			//	tool.LogInfo.Println(f)
			//}


			boolReplace, err := tool.Confirm("table [" + table + "] exists need replace?", 3)
			if err != nil{
				return err
			}
			// replace old  table
			if boolReplace{
				err := createTable(database, createTableSql); if err != nil{
					return err
				}
			} else {
				// not replace old table
				tool.LogWarning.Println("table [" + table + "] exists")
			}
		}
	}
	return nil
}


func createTable(database string, createSql []string) error {
	for _, execSql := range createSql{
		if execSql == ""{
			continue
		}
		tool.LogInfo.Println(execSql)
		_, err := DbPool.Exec(execSql); if err != nil{
			return err
		}
	}
	return nil
}