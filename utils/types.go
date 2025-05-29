package utils

import "os"

type GetEnvFunc func(key string) string
type MkdirFunc func(path string, perm os.FileMode) error
type StatFunc func(name string) (os.FileInfo, error)
type ReadDirFunc func(name string) ([]os.DirEntry, error)
type RemoveFunc func(name string) error
