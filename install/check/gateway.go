package check

import (
	"fmt"
	"gatekeeper/install/tool"
	"io"
	"net/http"
	"os"
)


var (
	GitHubUrl  		string = "https://github.com/didi/Gatekeeper/archive/refs/heads/master.zip"
	GitClone		string = "https://github.com/didi/Gatekeeper.git"
	GatekeeperZip 	string = "gatekeeper.zip"
	InstallDestDir	string = tool.GetCurrentPath()
	InstallName		string = "Gatekeeper"
	CmdRun			string = fmt.Sprintf("cd %s/%s; %s run main.go run -c %s/%s/conf/dev/ -p control",
		 InstallDestDir, InstallName, GoPath, InstallDestDir, InstallName)
)


func InitGateway() error {
	InstallDestDir, err:= tool.Input("please enter install dir (default:" + InstallDestDir + "):", InstallDestDir)
	if err != nil{
		return err
	}
	tool.LogInfo.Println("install path: " + InstallDestDir)

	err = gitClone(); if err != nil{
		tool.LogWarning.Println(err)
		err = download(); if err != nil{
			return err
		}
	}
	return nil
}


func download() error{
	_, err = os.Stat(InstallDestDir)
	if err != nil {
		if os.IsNotExist(err) {
			tool.LogInfo.Println("install dir not exists")
			tool.LogInfo.Println("create install dir :" + InstallDestDir)
			err = os.Mkdir(InstallDestDir, 0666)
			if err != nil {
				return err
			}
		}
	}

	tool.LogInfo.Println("download gatekeeper service from: " + GitHubUrl)

	res, err := http.Get(GitHubUrl)
	tool.LogInfo.Println("download gatekeeper service start")

	if err != nil {
		//panic(err)
		tool.LogInfo.Println("download gatekeeper service error")
		return err
	}
	f, err := os.Create(InstallDestDir + "/" + GatekeeperZip)
	if err != nil {
		tool.LogInfo.Println("download gatekeeper service error")
		//panic(err)
		return err
	}

	_, err = io.Copy(f, res.Body)
	if err != nil{
		tool.LogInfo.Println("download gatekeeper service error")
		//return err
	}
	tool.LogInfo.Println("download gatekeeper service end")

	tool.LogInfo.Println("unpack gatekeeper service start")
	err = unzip()
	if err != nil{
		return err
	}
	tool.LogInfo.Println("unpack gatekeeper service end")

	return nil
}


func gitClone() error{
	cmdClone := fmt.Sprintf("git clone %s %s/%s", GitClone, InstallDestDir, InstallName)
	tool.LogInfo.Println(cmdClone)
	err := tool.RunCmd(cmdClone)
	if err != nil{
		return err
	}
	return nil
}


func unzip() error {
	err = tool.Unzip(InstallDestDir + "/" + GatekeeperZip, InstallDestDir)
	if err != nil{
		return err
	}
	err = os.Rename(InstallDestDir + "/Gatekeeper-master", InstallDestDir + "/" + InstallName)
	if err != nil{
		return err
	}
	return nil
}


func RunGateKeeper() error{
	boolRunGatekeeper, err := tool.Confirm("run gatekeeper?", 2)
	if err != nil{
		return err
	}
	if boolRunGatekeeper {
		tool.LogInfo.Println(CmdRun)
		err := tool.RunCmd(CmdRun)
		if err != nil{
			return err
		}
	}
	return nil
}