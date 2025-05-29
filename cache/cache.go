package cache

import (
	"os"
	"path/filepath"
	"sort"

	"example.com/termquery/utils"
	"github.com/google/uuid"
)

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

		for _, element := range fileList {
			queue.Enqueue(element.Name())
		}

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
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
