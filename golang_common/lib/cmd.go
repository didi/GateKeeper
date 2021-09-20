package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var confPath string
var panelType string
var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run a gatekeeper application",
	Long:  `Run a gatekeeper application by parameter`,
	Args: func(cmd *cobra.Command, args []string) error {
		panelType = cmd.Flag("panel_type").Value.String()
		confPath = cmd.Flag("conf_path").Value.String()
		if ok, _ := PathExists(confPath); !ok {
			return errors.New("conf_path is not a real dir")
		}
		if !InArrayString(panelType, []string{"proxy", "control"}) {
			return errors.New("panel_type error，Choose one from proxy, control")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func CmdExecute() error {
	cmdRun.Flags().StringVarP(&panelType, "panel_type", "p", "", "Set panel type(control/proxy)")
	cmdRun.Flags().StringVarP(&confPath, "conf_path", "c", "", "Set configuration path(./conf/dev/)")
	cmdRun.MarkFlagRequired("panel_type")
	cmdRun.MarkFlagRequired("conf_path")
	var rootCmd = &cobra.Command{
		Use:               "",
		Short:             "Gatekeeper command manager",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true, DisableNoDescFlag: true, DisableDescriptions: true},
	}
	rootCmd.AddCommand(cmdRun)
	gin.SetMode(gin.ReleaseMode)
	return rootCmd.Execute()
}

func GetCmdConfPath() string {
	return confPath
}

func GetCmdPanelType() string {
	return panelType
}
