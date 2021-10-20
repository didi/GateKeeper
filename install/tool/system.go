package tool

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	SystemType  string
	CommandType string
	CommandArgs string
	GateKeeperPath	string = gatekeeperPath()
)


func InitSystem()  {
	SystemType = runtime.GOOS
	if SystemType == "windows"{
		CommandType = "cmd"
		CommandArgs = "/C"
	} else {
		CommandType = "sh"
		CommandArgs = "-c"
	}
}

func Cmd(command string) (string, error){
	cmd := exec.Command(CommandType, CommandArgs, command)
	str, err := cmd.Output()
	LogInfo.Println(string(str))
	return string(str), err
}


func RunCmd(command string) error{
	cmd := exec.Command(CommandType, CommandArgs, command)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		LogInfo.Println(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}


func GetCurrentPath() string{
	path, _ := os.Getwd()
	return strings.Replace(path, "\\", "/", -1)
}


func gatekeeperPath() string{
	path := GetCurrentPath()
	pathArr := strings.Split(path, "/")
	index := len(pathArr)
	pathArr = pathArr[0:index-1]
	path = strings.Join(pathArr, "/")
	return path
}