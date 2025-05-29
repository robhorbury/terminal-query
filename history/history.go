package history

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"

	"example.com/termquery/constants"
	"example.com/termquery/utils"
)

func GetHomeDir(envFunc utils.GetEnvFunc) (string, error) {
	envVar := envFunc("HOME")
	if envVar == "" {
		return "", fmt.Errorf("$HOME environment variable must be set")
	} else {
		return envVar, nil
	}
}

func GetCacheDir(homeDir string, envFunc utils.GetEnvFunc) string {
	envVar := envFunc("XDG_CACHE_HOME")
	if envVar == "" {
		return filepath.Join(homeDir, constants.DefaultXDGCacheDirectory, constants.DefaultApplicationCacheDirectory)
	} else {
		return filepath.Join(envVar, constants.DefaultApplicationCacheDirectory)
	}
}

func GetForceUseNeovim(envFunc utils.GetEnvFunc) bool {
	envVar := envFunc("TERMQUERY_FORCE_USE_NEOVIM")
	if envVar == "" {
		return false
	} else if envVar == "true" {
		return true
	} else {
		return false
	}
}

func GetEditor(forceUseNvimFlag bool, envFunc utils.GetEnvFunc) string {
	if forceUseNvimFlag == true {
		return "nvim"
	}
	envVar := envFunc("EDITOR")
	if envVar == "" {
		return "vi"
	} else {
		return envVar
	}
}

func getMaxNumberOfHistoricalQueries(envFunc utils.GetEnvFunc) int16 {
	envVar := envFunc("TERMQUERY_HISTORICAL_QUERY_LIMIT")
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

func InitCache(cachePath string, stat utils.StatFunc, mkdir utils.MkdirFunc) error {
	if !utils.FolderExists(cachePath, stat) {
		error := mkdir(cachePath, os.ModePerm)
		return error
	} else {
		return nil
	}
}

func CreateFileQueue(cachePath string, maxNumberQueries int16, readDirFunc utils.ReadDirFunc, removeFunc utils.RemoveFunc) (*utils.FileQueue, error) {
	fileList, err := readDirFunc(cachePath)
	if err != nil {
		return nil, err
	} else {

		sort.Slice(fileList, func(i, j int) bool {
			file1Stat, err := fileList[i].Info()
			if err != nil {
				panic(err)
			}
			file2Stat, err := fileList[j].Info()
			if err != nil {
				panic(err)
			}
			return file1Stat.ModTime().Unix() < file2Stat.ModTime().Unix()
		})

		queue := utils.NewFileQueue()

		for _, element := range fileList {
			queue.Enqueue(element.Name())
		}

		for queue.Length >= int(maxNumberQueries) {
			queue.RemoveAndDeque(cachePath, removeFunc)
		}

		return queue, nil
	}
}

func EditFile(editor string, filePath string, commandFunc utils.CommandFunc) {
	// Open editor in blocking mode
	cmd := commandFunc(editor, filePath)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
