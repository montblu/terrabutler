package main

import (
	"os"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

// Checks the requirements before running the application.
func check_requirement() {

	// Sync logger
	defer logger.Sync()

	// Loading koanf instance
	var k = koanf.New(".")

	//Getting the environment variables
	k.Load(env.Provider(".", env.Opt{Prefix: "TERRABUTLER_"}), nil)
	root := k.String("TERRABUTLER_ROOT")
	isEnabled := k.Bool("TERRABUTLER_ENABLE")
	settingsFile := root + "/configs/settings.yml"

	logger.Info("Environment Variables:", zap.String("TERRAFORM_ROOT", root), zap.Bool("TERRAFORM_ENABLED", isEnabled), zap.String("Settings Location", settingsFile))

	if !isEnabled {
		logger.Error("Terrabutler is not currently enabled on this folder. Please set 'TERRABUTLER_ENABLE' in your environment to true to enable it.")
		os.Exit(1)
	}
	if root == "" {
		logger.Error("Terrabutler can't determine the root folder of your project or it doesn't exist. Please set 'TERRABUTLER_ROOT' in your environment pointing to the root folder of your project.")
		os.Exit(1)
	}
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		logger.Error("Terrabutler can't find you settings file. Please create a 'settings.yml' file inside the 'configs' folder.")
		os.Exit(1)
	}
}
