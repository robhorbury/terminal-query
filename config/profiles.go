package config

import (
	"fmt"
	"path"
	"strings"

	"example.com/termquery/constants"
)

func readProfile(params ConfigParams, profileName string) ([]string, error) {
	fileContents, err := params.ReadFileFunc(path.Join(params.ConfigPath, constants.ProfilesFileName))

	if !strings.Contains(string(fileContents), "["+profileName+"]") {
		return []string{""}, fmt.Errorf("Profile %s not in profiles", profileName)
	}
	if err != nil {
		return []string{""}, err
	}

	profileStrings := strings.Split(strings.Split(strings.Split(string(fileContents), "["+profileName+"]")[1], "[")[0], "\n")
	return profileStrings, nil
}

func extractSettings(profileStrings []string, settingsName string) (string, error) {
	for i := range profileStrings {
		if strings.Contains(profileStrings[i], settingsName) {
			settingSplit := strings.Split(profileStrings[i], ":")
			return settingSplit[1], nil
		}
	}
	return "", fmt.Errorf("No setting matching %s found", settingsName)
}

func GetHttpPath(params ConfigParams, profileName string) (string, error) {
	profileString, err := readProfile(params, profileName)
	if err != nil {
		return "", err
	}
	settingValue, err := extractSettings(profileString, "http_path")
	return settingValue, err
}

func GetToken(params ConfigParams, profileName string) (string, error) {
	profileString, err := readProfile(params, profileName)
	if err != nil {
		return "", err
	}
	settingValue, err := extractSettings(profileString, "access_token")
	return settingValue, err
}

func GetServerHostname(params ConfigParams, profileName string) (string, error) {
	profileString, err := readProfile(params, profileName)
	if err != nil {
		return "", err
	}
	settingValue, err := extractSettings(profileString, "server_hostname")
	return settingValue, err
}
