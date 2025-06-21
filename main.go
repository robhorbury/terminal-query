package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"example.com/termquery/cache"
	"example.com/termquery/config"
	"example.com/termquery/logger"
	"example.com/termquery/sql"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type errMsg error

type model struct {
	spinner         spinner.Model
	userQuitting    bool
	channelQuitting bool
	spinnerChannel  chan bool
	err             error
}

func initialModel(channel chan bool) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:        s,
		spinnerChannel: channel}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.spinnerChannel:
		m.channelQuitting = true
		return m, tea.Quit
	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "c", "ctrl+c":
				m.userQuitting = true
				os.Exit(1)
				tea.ClearScreen()
				return m, tea.Quit
			default:
				return m, nil
			}

		case errMsg:
			m.err = msg
			return m, nil

		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Running Query. Press c to cancel\n\n", m.spinner.View())
	if m.userQuitting || m.channelQuitting {
		return str + "\n"
	}
	return str
}

func RunQueryFromFileWithChannel(
	cacheParams cache.CacheParams,
	fileName string,
	connection sql.Connection,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	rowChannel chan []map[string]string,
	colChannel chan []string,
	errorChannel chan error,
	spinnerChannel chan bool,
) {
	defer wg.Done()

	rows, columns, err := connection.RunQueryFromFile(filepath.Join(cacheParams.CachePath, fileName))
	spinnerChannel <- true
	rowChannel <- rows
	colChannel <- columns
	errorChannel <- err
}

type TestQuery struct {
	col float32
}

func main() {

	spinnerFinished := make(chan bool, 1)
	rowChan := make(chan []map[string]string, 1)
	colChan := make(chan []string, 1)
	errorChan := make(chan error, 1)

	logger.Init(logger.LoggerConfig{
		Level:  slog.LevelError,
		Format: logger.FormatJSON, // or logger.FormatText
	})
	logger := logger.Get()

	home, err := cache.GetHomeDir(os.Getenv, logger)
	if err != nil {
		panic(err)
	}

	configParams := config.ConfigParams{
		Logger:        logger,
		ConfigPath:    cache.GetConfigDir(home, os.Getenv, logger),
		ReadDirFunc:   os.ReadDir,
		MkdirFunc:     os.MkdirAll,
		StatFunc:      os.Stat,
		WriteFileFunc: os.WriteFile,
		ReadFileFunc:  os.ReadFile,
	}

	config.InitConfig(configParams)

	cacheParams := cache.CacheParams{
		Logger:           logger,
		CachePath:        cache.GetCacheDir(home, os.Getenv, logger),
		MaxNumberQueries: config.GetMaxNumberHistoricalQueries(configParams),
		Editor:           cache.GetEditor(config.GetForceUseNeovim(configParams), os.Getenv, logger),
		RemoveFunc:       os.Remove,
		CommandFunc:      RealCommandFactory,
		ReadDirFunc:      os.ReadDir,
		MkdirFunc:        os.MkdirAll,
		StatFunc:         os.Stat,
	}

	token, err := config.GetToken(configParams, config.GetDefaultProfile(configParams))
	if err != nil {
		panic(err)
	}
	httpPath, err := config.GetHttpPath(configParams, config.GetDefaultProfile(configParams))
	if err != nil {
		panic(err)
	}
	ServerHostname, err := config.GetServerHostname(configParams, config.GetDefaultProfile(configParams))
	if err != nil {
		panic(err)
	}

	connection := sql.DatabricksConnection{
		AccessToken:    token,
		HttpPath:       httpPath,
		ServerHostname: ServerHostname,
		Logger:         logger,
	}

	cache.InitCache(cacheParams)
	queue, _ := cache.CreateFileQueue(cacheParams)

	// file_name := cache.CreateAndEnque(queue, cacheParams, cache.EditFile)
	file_name := cache.EditMostRecentFile(queue, cacheParams, cache.EditFile)
	var wg sync.WaitGroup
	wg.Add(1)
	go RunQueryFromFileWithChannel(cacheParams, file_name, connection, &wg, logger, rowChan, colChan, errorChan, spinnerFinished)
	model := initialModel(spinnerFinished)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}

	wg.Wait()
	rows := <-rowChan
	columns := <-colChan
	err = <-errorChan

	// sql.PrintRowsAsTableBasic(os.Stdout, rows)
	sql.PrintRowsAsTableTea(rows, columns)
}
