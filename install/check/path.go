package check

import (
	"errors"
	"gatekeeper/install/tool"
	"regexp"
	"strconv"
	"strings"
)

var (
	UseVersion = "1.12.0"
	GoPath string = "go"
)

func InitGo() error{
	err := checkVersion()
	if err != nil{
		tool.LogWarning.Println(err)
		err = enterGoPath()
		if err != nil{
			return err
		}

		err = checkVersion()
		if err != nil{
			return err
		}
	}

	return nil
}


func checkVersion() error{

	str, err := tool.Cmd(GoPath + " version")
	if err != nil{
		return errors.New("go command not found")
	}
	versionPre := regexp.MustCompile(`[0-9]\.[0-9]{1,2}\.[0-9]*`)
	match := versionPre.FindString(string(str))
	goVersion, _ := strconv.Atoi(strings.Replace(match, ".", "", -1))
	useVersion, _ := strconv.Atoi(strings.Replace(UseVersion, ".", "", -1))
	if goVersion < useVersion {
		return errors.New("gatekeeper use go version must be >= 1.12.0 please check")
	}
	return nil
}


func enterGoPath() error{
	goPath, err := tool.Input("please enter go path (/use/bin/go):", "")
	if err != nil{
		tool.LogError.Println("inner error")
		return err
	}
	if goPath == ""{
		return errors.New("no go use")
	}
	GoPath = goPath
	return nil
}