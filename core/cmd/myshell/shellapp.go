package main

import (
	"context"

	"github.com/duxlabs-dev/dux-go/shellapp"
	"github.com/duxlabs-dev/dux-shell-starter/core/internal/backend"
	"github.com/duxlabs-dev/dux-shell-starter/core/internal/qtbaseembed"
	"github.com/duxlabs-dev/dux-shell-starter/core/internal/shellembed"
)

const (
	appID     = "myshell"
	envPrefix = "MYSHELL"
	desktopID = "com.example.myshell"
)

var shellApp = shellapp.New(shellapp.Config{
	ID:        appID,
	EnvPrefix: envPrefix,
	QSAppID:   desktopID,
	Version:   Version,
	Embedded:  embeddedShell{},
	QtBase:    qtbaseembed.Extract,
	Boot:      bootBackend,
})

func bootBackend(ctx context.Context) (shellapp.Backend, error) {
	return backend.Start(ctx, appID, Version)
}

type embeddedShell struct{}

func (embeddedShell) Available() bool { return shellembed.Available() }

func (embeddedShell) Extract(baseDir string) (string, error) { return shellembed.Extract(baseDir) }

func (embeddedShell) Prune(baseDir, keep string) { shellembed.Prune(baseDir, keep) }
