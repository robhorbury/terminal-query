package utils

import (
	"os"
	"os/exec"
)

type MkdirFunc func(path string, perm os.FileMode) error
type StatFunc func(name string) (os.FileInfo, error)
type GetEnvFunc func(key string) string
type ReadDirFunc func(name string) ([]os.DirEntry, error)
type RemoveFunc func(name string) error
type CommandFunc func(name string, arg ...string) *exec.Cmd
