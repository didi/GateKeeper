package public

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/golang_common/log"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v8"
	"time"
)

//公共handle
var (
	ConfHandler      *viper.Viper
	ConfPath         string
	AuthConf         *AuthConfig
	Uptime           time.Time
	ValidatorHandler = validator.New(&validator.Config{TagName: "validate", FieldNameTag: "json"})
	SysLogger        *log.Logger
	CheckLogger      *log.Logger
	StatLogger       *log.Logger
	TraceLoggerOn    bool

	TimeLocation   *time.Location
	FullTimeFormat = "2006-01-02 15:04:05"
	DayFormat      = "2006-01-02"
	TimeFormat     = "15:04:05"
)

//ContextKey context key type
type ContextKey string

//InitConf 初始化initconf
func InitConf() {
	Uptime = time.Now()
	SysLogger = log.NewLogger()
	log.SetupLogInstanceWithConf(log.LogConfig{
		Level: lib.GetStringConf("base.syslog.log_level"),
		FW: log.ConfFileWriter{
			On:              lib.GetBoolConf("base.syslog.file_writer.on"),
			LogPath:         lib.GetStringConf("base.syslog.file_writer.log_path"),
			RotateLogPath:   lib.GetStringConf("base.syslog.file_writer.rotate_log_path"),
			WfLogPath:       lib.GetStringConf("base.syslog.file_writer.wf_log_path"),
			RotateWfLogPath: lib.GetStringConf("base.syslog.file_writer.rotate_wf_log_path"),
		},
		CW: log.ConfConsoleWriter{
			On:    lib.GetBoolConf("base.syslog.console_writer.on"),
			Color: lib.GetBoolConf("base.syslog.console_writer.color"),
		},
	}, SysLogger)

	CheckLogger = log.NewLogger()
	log.SetupLogInstanceWithConf(log.LogConfig{
		Level: lib.GetStringConf("base.checklog.log_level"),
		FW: log.ConfFileWriter{
			On:              lib.GetBoolConf("base.checklog.file_writer.on"),
			LogPath:         lib.GetStringConf("base.checklog.file_writer.log_path"),
			RotateLogPath:   lib.GetStringConf("base.checklog.file_writer.rotate_log_path"),
			WfLogPath:       lib.GetStringConf("base.checklog.file_writer.wf_log_path"),
			RotateWfLogPath: lib.GetStringConf("base.checklog.file_writer.rotate_wf_log_path"),
		},
		CW: log.ConfConsoleWriter{
			On:    lib.GetBoolConf("base.checklog.console_writer.on"),
			Color: lib.GetBoolConf("base.checklog.console_writer.color"),
		},
	}, CheckLogger)

	StatLogger = log.NewLogger()
	log.SetupLogInstanceWithConf(log.LogConfig{
		Level: lib.GetStringConf("base.statlog.log_level"),
		FW: log.ConfFileWriter{
			On:              lib.GetBoolConf("base.statlog.file_writer.on"),
			LogPath:         lib.GetStringConf("base.statlog.file_writer.log_path"),
			RotateLogPath:   lib.GetStringConf("base.statlog.file_writer.rotate_log_path"),
			WfLogPath:       lib.GetStringConf("base.statlog.file_writer.wf_log_path"),
			RotateWfLogPath: lib.GetStringConf("base.statlog.file_writer.rotate_wf_log_path"),
		},
		CW: log.ConfConsoleWriter{
			On:    lib.GetBoolConf("base.statlog.console_writer.on"),
			Color: lib.GetBoolConf("base.statlog.console_writer.color"),
		},
	}, StatLogger)

	AuthConf = &AuthConfig{}
	if err := lib.ParseLocalConfig("admin.toml", AuthConf); err != nil {
		log.Fatal("conf init error: ", err)
	}

	tl, terr := time.LoadLocation("Asia/Chongqing")
	if terr != nil {
		log.Fatal("conf init error: ", terr)
	}
	TimeLocation = tl
	FlowLimiterHandler = NewFlowLimiter()
	FlowCounterHandler = NewFlowCounter()

	TraceLoggerOn = true
	if lib.GetStringConf("base.base.access_log") == "off" {
		TraceLoggerOn = false
	}
}

//AuthConfig 验证结构体
type AuthConfig struct {
	Base struct {
		AdminName     string `mapstructure:"admin_username"`
		AdminPassport string `mapstructure:"admin_passport"`
	} `mapstructure:"base"`
}

//IsProductEnv 生产环境
func IsProductEnv() bool {
	if lib.GetConfEnv() == "prod" {
		return true
	}
	return false
}
