package cache

import (
	"io"

	"example.com/termquery/utils"
	"log/slog"
)

type Command interface {
	Run() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

type CacheParams struct {
	Logger           *slog.Logger
	CachePath        string
	MaxNumberQueries int16
	Editor           string
	RemoveFunc       utils.RemoveFunc
	CommandFunc      CommandFunc
	ReadDirFunc      utils.ReadDirFunc
	MkdirFunc        utils.MkdirFunc
	StatFunc         utils.StatFunc
}

type CommandFunc func(name string, arg ...string) Command

type EditFileFunc func(filePath string, params CacheParams) error
