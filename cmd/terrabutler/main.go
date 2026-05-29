package main

import (
	"github.com/l58193/terrabutler/internal/cli"
	"github.com/l58193/terrabutler/internal/logger"

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

	err := cli.Run(version, fs)
	if err != nil {
		logger.Zap.Error(err.Error())
	}
}
