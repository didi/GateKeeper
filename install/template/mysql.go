package template

import (
	"bufio"
	"gatekeeper/install/tool"
	"os"
	"regexp"
	"strings"
)

var(
	TableSql map[string]string
	Tables map[string][]string
)

func InitSql() error {
	TableSql = make(map[string]string)
	Tables = make(map[string][]string)
	err := getCreateSql()
	if err != nil{
		return err
	}
	for tableName, sql := range TableSql {
		Tables[tableName] = strings.Split(sql, ";")
	}
	return nil
}


func getCreateSql() error{
	sqlFilePath := tool.GateKeeperPath + "/gatekeeper.sql"
	tool.LogInfo.Println("sql file path :" + sqlFilePath)
	tableName := ""
	tablePre := regexp.MustCompile(`gateway_[a-z_]*`)
	f, _ := os.Open(sqlFilePath)
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		lineSql, err := tool.ReadLine(r)
		if strings.Index(lineSql, "/*") < 0 && strings.Index(lineSql, "--") < 0 && lineSql != "" {
			if strings.Index(lineSql, "DROP TABLE IF EXISTS") >= 0 {
				tableName = tablePre.FindString(lineSql)
			}
			TableSql[tableName] += lineSql
		}
		if err != nil {
			break
		}
	}
	return nil
}