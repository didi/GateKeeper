package main

import (
	"gatekeeper/install/check"
	"gatekeeper/install/tool"
	"os"
)

func main() {
	var (
		err error
	)

	tool.InitSystem()

	err = check.InitRedis(); if err != nil{
		tool.LogError.Println("connect redis error")
		tool.LogError.Println(err)
		os.Exit(-1)
	}

	err = check.InitDb(); if err != nil{
		tool.LogError.Println(err)
		os.Exit(-1)
	}

	err = check.InitConf(); if err != nil{
		tool.LogError.Println(err)
		os.Exit(-1)
	}

	err = check.InitGo(); if err != nil{
		tool.LogWarning.Println(err)
		os.Exit(-1)
	}

	err = check.RunGateKeeper(); if err != nil{
		tool.LogError.Println(err)
		os.Exit(-1)
	}

}

