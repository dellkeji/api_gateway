package main

// NOTE: panic-recovery is cost resource, so use in exception function

import (
	"os"

	cmd "apigw_golang/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
