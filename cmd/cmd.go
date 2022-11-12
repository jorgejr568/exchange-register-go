package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "exchange-register-go",
	Short: "Exchange Register Go",
}

func Execute() error {
	return rootCmd.Execute()
}
