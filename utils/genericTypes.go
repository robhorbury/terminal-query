package utils

import (
	"os"
	"os/exec"

	"log/slog"
)

type HistoryParams struct {
	Logger           *slog.Logger
	CachePath        string
	MaxNumberQueries int16
	Editor           string
	RemoveFunc       RemoveFunc
	CommandFunc      CommandFunc
	ReadDirFunc      ReadDirFunc
	MkdirFunc        MkdirFunc
	StatFunc         StatFunc
}

type MkdirFunc func(path string, perm os.FileMode) error
type StatFunc func(name string) (os.FileInfo, error)
type GetEnvFunc func(key string) string
type ReadDirFunc func(name string) ([]os.DirEntry, error)
type RemoveFunc func(name string) error
type CommandFunc func(name string, arg ...string) *exec.Cmd

type EditFileFunc func(filePath string, params HistoryParams) error
