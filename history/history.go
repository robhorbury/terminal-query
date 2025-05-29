package history

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"example.com/termquery/constants"
	"example.com/termquery/utils"
	"github.com/google/uuid"
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

func InitCache(params utils.HistoryParams) error {
	if !utils.FolderExists(params.CachePath, params.StatFunc) {
		error := params.MkdirFunc(params.CachePath, os.ModePerm)
		return error
	} else {
		return nil
	}
}

func CreateFileQueue(params utils.HistoryParams) (*utils.FileQueue, error) {
	fileList, err := params.ReadDirFunc(params.CachePath)
	if err != nil {
		return nil, err
	} else {

		sort.Slice(fileList, func(i, j int) bool {
			file1Stat, err := fileList[i].Info()
			if err != nil {
				params.Logger.Error("Cound not sort file", "Index", i)
				panic(err)
			}
			file2Stat, err := fileList[j].Info()
			if err != nil {
				params.Logger.Error("Cound not sort file", "Index", j)
				panic(err)
			}
			return file1Stat.ModTime().Unix() < file2Stat.ModTime().Unix()
		})

		queue := utils.NewFileQueue()
		fmt.Println("HERE A: ")
		fmt.Println(fileList)

		for i, element := range fileList {
			fmt.Println("HERE: ", i)
			queue.Enqueue(element.Name())
		}
		fmt.Println("QUEUE LENGTH: ", queue.Length)

		for queue.Length > int(params.MaxNumberQueries) {
			queue.RemoveAndDeque(params.CachePath, params.RemoveFunc)
		}

		return queue, nil
	}
}

func CreateAndEnque(queue *utils.FileQueue, params utils.HistoryParams, editFunc utils.EditFileFunc) string {
	params.Logger.Debug("VAR:", "queue.length", queue.Length)
	fileName := uuid.New().String() + ".sql"

	if queue.Length < int(params.MaxNumberQueries) {
		params.Logger.Debug("Enqueue a new file")

		queue.Enqueue(fileName)
		editFunc(fileName, params)
	} else {
		params.Logger.Debug("Replace a file")
		err := queue.RemoveAndDeque(params.CachePath, params.RemoveFunc)
		if err != nil {
			params.Logger.Error("Could not remve and deque")
			panic(err)
		}
		queue.Enqueue(fileName)
		editFunc(fileName, params)
	}
	queue, err := CreateFileQueue(params)
	if err != nil {

		params.Logger.Error("Could not create file queue")
		panic(err)
	}

	return fileName
}

func EditFile(fileName string, params utils.HistoryParams) error {
	// Open editor in blocking mode
	cmd := params.CommandFunc(params.Editor, filepath.Join(params.CachePath, fileName))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
