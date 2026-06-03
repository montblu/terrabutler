package main

import (
	"github.com/montblu/terrabutler/internal/cli"
	"github.com/montblu/terrabutler/internal/logger"

	"github.com/spf13/afero"
)

func main() {

	// Using Real FileSystem
	fs := afero.NewOsFs()

	version := "v3.0.0"

	err := cli.Run(version, fs)

	if err != nil {
		logger.Zap.Error(err.Error())
	}
}
