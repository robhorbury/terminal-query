package cache

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"example.com/termquery/constants"
	"example.com/termquery/utils"
)

func GetHomeDir(envFunc utils.GetEnvFunc, logger *slog.Logger) (string, error) {
	envVar := envFunc("HOME")
	logger.Debug("ENV_VAR:", "HOME", envVar)
	if envVar == "" {
		return "", fmt.Errorf("$HOME environment variable must be set")
	} else {
		return envVar, nil
	}
}

func GetConfigDir(homeDir string, envFunc utils.GetEnvFunc, logger *slog.Logger) string {

	envVar := envFunc("XDG_CONFIG_HOME")
	logger.Debug("ENV_VAR:", "XDG_CONFIG_HOME", envVar)
	if envVar == "" {
		return filepath.Join(homeDir, constants.DefaultXDGConfigDirectiory, constants.DefaultApplicationConfigDirectory)

	} else {
		return filepath.Join(envVar, constants.DefaultApplicationConfigDirectory)
	}
}

func GetCacheDir(homeDir string, envFunc utils.GetEnvFunc, logger *slog.Logger) string {

	envVar := envFunc("XDG_CACHE_HOME")
	logger.Debug("ENV_VAR:", "XDG_CACHE_HOME", envVar)
	if envVar == "" {
		return filepath.Join(homeDir, constants.DefaultXDGCacheDirectory, constants.DefaultApplicationCacheDirectory)
	} else {
		return filepath.Join(envVar, constants.DefaultApplicationCacheDirectory)
	}
}

func GetEditor(forceUseNvimFlag bool, envFunc utils.GetEnvFunc, logger *slog.Logger) string {
	if forceUseNvimFlag == true {
		return "nvim"
	}
	envVar := envFunc("EDITOR")
	logger.Debug("ENV_VAR:", "EDITOR", envVar)
	if envVar == "" {
		return "vi"
	} else {
		return envVar
	}
}
