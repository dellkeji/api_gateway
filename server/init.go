package server

import (
	"github.com/spf13/cobra"

	router "apigw_golang/router"
)

// ServerStartCmd : start server
func ServerStartCmd() *cobra.Command {
	r := &router.Router{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start a api gateway server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.Init(ProxyHandler)
		},
	}

	return cmd
}
