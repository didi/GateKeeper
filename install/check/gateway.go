package check

import (
	"fmt"
	"gatekeeper/install/tool"
)


var (
	GateKeeperPath	string = tool.GateKeeperPath
	CmdRun			string = "cd %s && %s run main.go run -c %s/conf/dev/ -p control"
)


func RunGateKeeper() error{
	boolRunGatekeeper, err := tool.Confirm("run gatekeeper?", 2)
	if err != nil{
		return err
	}
	CmdRun := fmt.Sprintf(CmdRun, GateKeeperPath, GoPath, GateKeeperPath)
	if boolRunGatekeeper {
		tool.LogInfo.Println(CmdRun)
		err := tool.RunCmd(CmdRun)
		if err != nil{
			return err
		}
	}
	return nil
}