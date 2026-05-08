package main

import (
	"terrabutler/internal/cli"
	"terrabutler/internal/logger"

	"github.com/spf13/afero"
)

func main() {

	// Using Real FileSystem
	fs := afero.NewOsFs()

	version := "v3.0.1"

	err := cli.Run(version, fs)

	if err != nil {
		logger.Zap.Error(err.Error())
	}
}
