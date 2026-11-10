package main

import (
	"os"

	"github.com/duxlabs-dev/dux-go/log"
	"github.com/duxlabs-dev/dux-go/qtbase"
	"github.com/duxlabs-dev/dux-shell-starter/core/internal/qtbaseembed"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"
)

func main() {
	log.SetEnvPrefix(envPrefix)
	qtbase.MaybeRunAndExit(qtbaseembed.Extract)
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
