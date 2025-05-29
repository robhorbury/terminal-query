package cache

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"

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

func GetCacheDir(homeDir string, envFunc utils.GetEnvFunc, logger *slog.Logger) string {

	envVar := envFunc("XDG_CACHE_HOME")
	logger.Debug("ENV_VAR:", "XDG_CACHE_HOME", envVar)
	if envVar == "" {
		return filepath.Join(homeDir, constants.DefaultXDGCacheDirectory, constants.DefaultApplicationCacheDirectory)
	} else {
		return filepath.Join(envVar, constants.DefaultApplicationCacheDirectory)
	}
}

func GetForceUseNeovim(envFunc utils.GetEnvFunc, logger *slog.Logger) bool {
	envVar := envFunc("TERMQUERY_FORCE_USE_NEOVIM")
	logger.Debug("ENV_VAR:", "TERMQUERY_FORCE_USE_NEOVIM", envVar)
	if envVar == "" {
		return false
	} else if envVar == "true" {
		return true
	} else {
		return false
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

func GetMaxNumberOfHistoricalQueries(envFunc utils.GetEnvFunc, logger *slog.Logger) int16 {
	envVar := envFunc("TERMQUERY_HISTORICAL_QUERY_LIMIT")
	logger.Debug("ENV_VAR:", "TERMQUERY_HISTORICAL_QUERY_LIMIT", envVar)
	if envVar == "" {
		return constants.DefaultMaxNumberOfHistoricalQueries
	}
	variableAsInt, err := strconv.Atoi(envVar)
	if err != nil {
		return constants.DefaultMaxNumberOfHistoricalQueries
	} else if variableAsInt <= 0 || variableAsInt >= 32767 {
		return constants.DefaultMaxNumberOfHistoricalQueries
	} else {
		return int16(variableAsInt)
	}

}
