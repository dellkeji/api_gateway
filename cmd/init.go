package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	conf "apigw_golang/configure"
	server "apigw_golang/server"
)

var configFile string

// NewRootCmd creates a new instance of the root command
func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "apigw",
		Short: "API Gateway is an API Gateway",
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "default.yaml", "configure file for apigw")

	//configure
	readConfigureFile(configFile)

	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(server.ServerStartCmd())

	return cmd
}

// readConfigureFile : read configire content
func readConfigureFile(configFile string) {
	err := conf.InitConfig(configFile)
	if err != nil {
		fmt.Println("Load configure file error, detail: ", err)
		os.Exit(-1)
	}
}

// NewVersionCmd : print current version
func NewVersionCmd() *cobra.Command {
	version := conf.GlobalConfigurations.Version
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version information",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("apigw version is %s\n", version)
		},
	}
}
