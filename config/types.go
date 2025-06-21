package config

import (
	"log/slog"

	"example.com/termquery/utils"
)

type ConfigParams struct {
	Logger        *slog.Logger
	ConfigPath    string
	ReadDirFunc   utils.ReadDirFunc
	MkdirFunc     utils.MkdirFunc
	StatFunc      utils.StatFunc
	WriteFileFunc utils.WriteFileFunc
	ReadFileFunc  utils.ReadFileFunc
}
