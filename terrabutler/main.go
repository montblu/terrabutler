package main

import (
	"terrabutler/internal/cli"
	"terrabutler/internal/logger"

	"go.uber.org/zap"
)

func main() {

	version := "v3.0.1"

	err := cli.Run(version)

	if err != nil {
		logger.Zap.Error("An error has Occured: ", zap.Error(err))
	}
}
