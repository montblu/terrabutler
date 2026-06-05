package main

import (
	"github.com/montblu/terrabutler/internal/cli"
	"github.com/montblu/terrabutler/internal/logger"

	"github.com/spf13/afero"
)

var (
	appName = "terrabutler"
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fs := afero.NewOsFs()

	err := cli.Run(appName, version, commit, date, fs)
	if err != nil {
		logger.Zap.Error(err.Error())
	}
}
