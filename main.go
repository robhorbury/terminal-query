package main

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"example.com/termquery/cache"
	"example.com/termquery/logger"
	"example.com/termquery/sql"
)

type RealCommand struct {
	cmd *exec.Cmd
}

func (r *RealCommand) Run() error {
	return r.cmd.Run()
}
func (r *RealCommand) SetStdin(in io.Reader) {
	r.cmd.Stdin = in
}

func (r *RealCommand) SetStdout(out io.Writer) {
	r.cmd.Stdout = out
}

func (r *RealCommand) SetStderr(err io.Writer) {
	r.cmd.Stderr = err
}

func RealCommandFactory(name string, args ...string) cache.Command {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &RealCommand{cmd: cmd}
}

type TestQuery struct {
	col float32
}

func main() {

	logger.Init(logger.LoggerConfig{
		Level:  slog.LevelDebug,
		Format: logger.FormatJSON, // or logger.FormatText
	})
	logger := logger.Get()

	home, err := cache.GetHomeDir(os.Getenv, logger)
	if err != nil {
		panic(err)
	}

	cacheParams := cache.CacheParams{
		Logger:           logger,
		CachePath:        cache.GetCacheDir(home, os.Getenv, logger),
		MaxNumberQueries: cache.GetMaxNumberOfHistoricalQueries(os.Getenv, logger),
		Editor:           cache.GetEditor(cache.GetForceUseNeovim(os.Getenv, logger), os.Getenv, logger),
		RemoveFunc:       os.Remove,
		CommandFunc:      RealCommandFactory,
		ReadDirFunc:      os.ReadDir,
		MkdirFunc:        os.MkdirAll,
		StatFunc:         os.Stat,
	}

	connection := sql.DatabricksConnection{
		AccessToken:    sql.GetDatabricksEnvVar("TERMQUERY_DATABRICKS_TOKEN", os.Getenv),
		HttpPath:       sql.GetDatabricksEnvVar("TERMQUERY_DATABRICKS_HTTP_PATH", os.Getenv),
		ServerHostname: sql.GetDatabricksEnvVar("TERMQUERY_DATABRICKS_HOST", os.Getenv),
		Logger:         logger,
	}

	cache.InitCache(cacheParams)
	queue, _ := cache.CreateFileQueue(cacheParams)

	// file_name := cache.CreateAndEnque(queue, cacheParams, cache.EditFile)
	file_name := cache.EditMostRecentFile(queue, cacheParams, cache.EditFile)
	rows, columns, err := connection.RunQueryFromFile(filepath.Join(cacheParams.CachePath, file_name))

	if err != nil {
		panic(err)
	}
	// sql.PrintRowsAsTableBasic(os.Stdout, rows)
	sql.PrintRowsAsTableTea(rows, columns)
}
