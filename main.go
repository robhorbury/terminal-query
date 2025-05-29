package main

import (
	"log/slog"
	"os"
	"os/exec"

	"example.com/termquery/history"
	"example.com/termquery/logger"
	"example.com/termquery/utils"
)

func main() {

	logger.Init(logger.LoggerConfig{
		Level:  slog.LevelDebug,
		Format: logger.FormatJSON, // or logger.FormatText
	})

	logger := logger.Get()

	home, err := history.GetHomeDir(os.Getenv, logger)
	if err != nil {
		panic(err)
	}
	cacheDir := history.GetCacheDir(home, os.Getenv, logger)

	historyParams := utils.HistoryParams{
		Logger:           logger,
		CachePath:        cacheDir,
		MaxNumberQueries: history.GetMaxNumberOfHistoricalQueries(os.Getenv, logger),
		Editor:           history.GetEditor(history.GetForceUseNeovim(os.Getenv, logger), os.Getenv, logger),
		RemoveFunc:       os.Remove,
		CommandFunc:      exec.Command,
		ReadDirFunc:      os.ReadDir,
		MkdirFunc:        os.MkdirAll,
		StatFunc:         os.Stat,
	}

	history.InitCache(historyParams)

	queue, _ := history.CreateFileQueue(historyParams)

	history.CreateAndEnque(queue, historyParams, history.EditFile)

}
