package tool

import (
	"log"
	"os"
)

var (
	LogInfo  	*log.Logger
	LogError 	*log.Logger
	LogWarning 	*log.Logger
)

func init() {
	LogInfo = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	LogError = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
	LogWarning = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime)
}
