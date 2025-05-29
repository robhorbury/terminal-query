package main

import (
	"example.com/termquery/history"
	"os"
)

func main() {
	home, err := history.GetHomeDir(os.Getenv)
	if err != nil {
		panic(err)
	}
	cacheDir := history.GetCacheDir(home, os.Getenv)
	history.InitCache(cacheDir, os.Stat, os.MkdirAll)

	history.CreateFileQueue(cacheDir, int16(10), os.ReadDir, os.Remove)

	history.EditFile("nvim", "~/temp.txt")
}
