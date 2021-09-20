package ddlog

import (
	"errors"
	"fmt"
	"strings"

	"github.com/didi/gatekeeper/golang_common/zerolog"
)

type FileConfig struct {
	//auto clear the log
	AutoClear bool

	//separate will  write wf log
	Separate bool

	//file directory
	FileDir string

	//file prefix name
	FilePrefix string

	//clear hours
	ClearHours int32

	//clear steps
	ClearStep int32

	//is Disable link
	DisableLink bool

	// log level
	Level string

	// logtype is file or stdout
	LogType string
}

func NewWriter(config *FileConfig) (zerolog.LevelWriter, error) {
	var writer zerolog.LevelWriter

	if err := checkConf(config); err != nil {
		return nil, err
	}

	if config.Separate {
		writer = zerolog.MultiLevelWriter(
			zerolog.NewAccessWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)),
			zerolog.NewWFFileWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)))

	} else {
		writer = zerolog.MultiLevelWriter(
			zerolog.NewFileWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)))
	}

	return writer, nil
}

func checkConf(config *FileConfig) error {
	if config.FileDir == "" || config.FilePrefix == "" || config.ClearHours <= 0 {
		errstr := fmt.Sprintf("log config validate failed， FileDir and FilePrefix should not be empty, ClearHours should more than 0, config is %+v", config)
		return errors.New(errstr)
	}
	return nil
}

func NewLoggerWithCfg(config *FileConfig) (*DiLogHandle, error) {
	var writer zerolog.LevelWriter

	if config.LogType == LogTypeStdout {
		writer = zerolog.MultiLevelWriter(zerolog.NewStdoutWriter())
	} else {
		if err := checkConf(config); err != nil {
			return nil, err
		}
		if config.Separate {
			writer = zerolog.MultiLevelWriter(
				zerolog.NewAccessWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)),
				zerolog.NewWFFileWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)))
		} else {
			writer = zerolog.MultiLevelWriter(
				zerolog.NewFileWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink)))
		}
	}

	diLogger := zerolog.New(writer)
	lev := strings.ToLower(config.Level)
	l, err := zerolog.ParseLevel(lev)
	if err != nil || lev == "" {
		l = zerolog.DebugLevel
	}
	diLogger = diLogger.Level(l)
	return &DiLogHandle{Logger: diLogger}, nil
}

func NewPubLogger(config *FileConfig) (*PubLog, error) {
	var log zerolog.Logger
	if config.LogType == LogTypeStdout {
		output := zerolog.NewStdoutWriter()
		output.FormatLevelFunc = func(i string) string {
			return ""
		}
		log = zerolog.New(output)
	} else {
		if err := checkConf(config); err != nil {
			return nil, err
		}
		output := zerolog.NewFileWriter(zerolog.SetAutoClear(config.AutoClear), zerolog.SetClearHours(config.ClearHours), zerolog.SetClearSteps(config.ClearStep), zerolog.SetFileDir(config.FileDir), zerolog.SetFilePrefix(config.FilePrefix), zerolog.SetDisableLink(config.DisableLink))

		output.FormatLevelFunc = func(i string) string {
			return ""
		}
		log = zerolog.New(output)
	}

	pLog := PubLog{
		Logger: log,
	}
	return &pLog, nil

}
