package config

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"example.com/termquery/constants"
	"example.com/termquery/utils"
)

func InitConfig(params ConfigParams) error {
	if !utils.FolderExists(params.ConfigPath, params.StatFunc) {
		error := params.MkdirFunc(params.ConfigPath, os.ModePerm)
		if error != nil {
			return error
		}
		error = createDefaultConfig(params)
		return error
	} else {
		if !utils.FileExists(path.Join(params.ConfigPath, constants.ConfigFileName), params.StatFunc) {
			createDefaultConfig(params)
		}
		return nil
	}
}

func createDefaultConfig(params ConfigParams) error {
	defaultConfig := fmt.Appendf([]byte(""), "max_number_historical_queries:%s\nforce_use_neovim:false\ndefault_profile:%s",
		strconv.FormatInt(int64(constants.DefaultMaxNumberOfHistoricalQueries), 10),
		constants.DefaultProfileName)
	err := params.WriteFileFunc(path.Join(params.ConfigPath, constants.ConfigFileName), defaultConfig, 0644)
	if err != nil {
		return err
	}
	err = params.WriteFileFunc(path.Join(params.ConfigPath, constants.ProfilesFileName), []byte(""), 0644)

	return err
}

func parseConfigFile(key string, params ConfigParams) (string, error) {
	fileContents, err := params.ReadFileFunc(path.Join(params.ConfigPath, constants.ConfigFileName))
	if err != nil {
		return "", err
	}
	fileContentsString := string(fileContents)
	fileContentsLines := strings.Split(fileContentsString, "\n")

	for i := range fileContentsLines {
		lineSplit := strings.Split(fileContentsLines[i], ":")
		if key == lineSplit[0] {
			return lineSplit[1], nil
		}
	}

	return "", fmt.Errorf("No key matches %s", key)
}

func GetMaxNumberHistoricalQueries(params ConfigParams) int16 {
	configValue, error := parseConfigFile("max_number_historical_queries", params)
	params.Logger.Debug("CONFIG:", "max_number_historical_queries", configValue)
	if error != nil {
		return constants.DefaultMaxNumberOfHistoricalQueries
	}

	if configValue == "" {
		return constants.DefaultMaxNumberOfHistoricalQueries
	}
	variableAsInt, err := strconv.Atoi(configValue)
	if err != nil {
		return constants.DefaultMaxNumberOfHistoricalQueries
	} else if variableAsInt <= 0 || variableAsInt >= 32767 {
		return constants.DefaultMaxNumberOfHistoricalQueries
	} else {
		return int16(variableAsInt)
	}
}

func GetForceUseNeovim(params ConfigParams) bool {
	configValue, error := parseConfigFile("force_use_neovim", params)
	params.Logger.Debug("CONFIG:", "force_use_neovim", configValue)
	if error != nil {
		return false
	}
	if configValue == "true" {
		return true
	}
	return false
}

func GetDefaultProfile(params ConfigParams) string {
	configValue, error := parseConfigFile("default_profile", params)
	params.Logger.Debug("CONFIG:", "default_profile", configValue)
	if error != nil {
		return constants.DefaultProfileName
	}
	return configValue
}
