package main

import (
	"os"
	"slices"

	"go.uber.org/zap"
)

func get_current_env() string {

	//Open site environment file
	env, err := os.ReadFile(paths["environment"])
	if err != nil {
		logger.Error("An error has occurred:", zap.Error(err))
	}
	logger.Info("Current environment is " + string(env))
	return string(env)
}

func is_protected_env(env string) bool {

	if slices.Contains(settings.Strings("environments.permanent"), env) {
		return true
	}
	return false
}
