package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           appID,
	Short:         "Dux shell starter",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(*cobra.Command, []string) {
		fmt.Printf("%s %s (commit %s, built %s)\n", appID, Version, Commit, BuildTime)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(shellApp.CustomConfigVar(), "config", "c", "", "Path to a UI config dir (containing shell.qml) to use instead of the embedded UI (env: "+envPrefix+"_SHELL_DIR)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(ipcCmd)
	rootCmd.AddCommand(shellApp.Commands()...)
}
