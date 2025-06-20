package sql

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

const textInputWidth = 50
const filterColumnWidth = 50
const tableColumnWidth = 20

// TODO: a better method of state managements
// TODO: I want to be able to see when filter is applied or not for example
// TODO: Change regex key from ? to \
// ─── Types & Constants ─────────────────────────────────────────────────────────

type Record = map[string]string

type viewState int

const (
	stateFiltering viewState = iota
	stateSelectVisibleColumns
	stateSelectFilterColumns
	stateNavigation
)

func newCustomDelegate() list.DefaultDelegate {
	return list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}
}

type columnItem struct {
	name     string
	selected bool
}

func (i columnItem) Title() string {
	mark := "[ ]"
	if i.selected {
		mark = "[x]"
	}
	return fmt.Sprintf("%s %s", mark, i.name)
}
func (columnItem) Description() string { return "" }
func (i columnItem) FilterValue() string {
	return i.name
}

// ─── Model ────────────────────────────────────────────────────────────────────

type model struct {
	state        viewState
	regexMode    bool // true when treating input as regex
	allCols      []string
	visibleCols  []string
	filterCols   []string
	allRows      []Record
	filteredRows []Record

	textInput   textinput.Model
	listVisible list.Model
	listFilter  list.Model
	table       table.Model
	keys        keyMap
	help        help.Model
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

// TODO: ADD MORE
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// TODO: Update keybindings!
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

// TODO: Add all keybinginds
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// NewModel constructs initial UI state.
func NewModel(data []map[string]string, cols []string) *model {
	// convert to Record
	rows := make([]Record, len(data))
	for i, r := range data {
		rows[i] = r
	}

	// prepare column‐picker lists
	// del := list.NewDefaultDelegate()
	del := newCustomDelegate()
	visibleItems := make([]list.Item, len(cols))
	filterItems := make([]list.Item, len(cols))
	for i, c := range cols {
		visibleItems[i] = columnItem{c, true}
		filterItems[i] = columnItem{c, true}
	}

	listVisible := list.New(visibleItems, &del, filterColumnWidth, len(visibleItems))
	listVisible.Title = "Visible Columns"
	listFilter := list.New(filterItems, &del, filterColumnWidth, len(filterItems))
	listFilter.Title = "Regex Columns"

	// text input for filtering
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = "> "
	ti.CharLimit = 128
	ti.Width = textInputWidth

	m := &model{
		state:        stateNavigation,
		regexMode:    false,
		allCols:      cols,
		visibleCols:  slices.Clone(cols),
		filterCols:   slices.Clone(cols),
		allRows:      rows,
		filteredRows: rows,
		textInput:    ti,
		listVisible:  listVisible,
		listFilter:   listFilter,
		keys:         keys,
		help:         help.New(),
	}

	m.textInput.Focus()
	m.buildTable(rows)
	return m
}

// PrintRowsAsTableTea starts the interactive TUI.
func PrintRowsAsTableTea(data []map[string]string, cols []string) {
	if len(data) == 0 {
		fmt.Println("No data to display.")
		return
	}
	p := tea.NewProgram(NewModel(data, cols))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

var (
	borderStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	headerStyle    = lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("250"))
	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
)

func (m *model) buildTable(rows []Record) {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 120
	}

	cols := make([]table.Column, len(m.visibleCols))
	for i, c := range m.visibleCols {
		cols[i] = table.NewColumn(c, c, tableColumnWidth)
	}

	tblRows := make([]table.Row, len(rows))
	for i, r := range rows {
		rd := table.RowData{}
		for _, col := range m.visibleCols {
			rd[col] = r[col]
		}
		tblRows[i] = table.NewRow(rd)
	}

	m.table = table.New(cols).
		WithRows(tblRows).
		WithMaxTotalWidth(w).
		WithPageSize(20).
		WithHorizontalFreezeColumnCount(1).
		WithMinimumHeight(10).
		Focused(true).
		WithBaseStyle(borderStyle).
		HeaderStyle(headerStyle).
		HighlightStyle(highlightStyle)
}

// columnList applies space/a/c/enter to a *list.Model and returns a
// slice of selected names only on Enter.
func columnList(l *list.Model, key string) []string {
	switch key {
	case " ":
		idx := l.Index()
		ci := l.Items()[idx].(columnItem)
		ci.selected = !ci.selected
		l.SetItem(idx, ci)

	case "a":
		for i := range len(l.Items()) {
			ci := l.Items()[i].(columnItem)
			ci.selected = true
			l.SetItem(i, ci)
		}

	case "c":
		for i := range len(l.Items()) {
			ci := l.Items()[i].(columnItem)
			ci.selected = false
			l.SetItem(i, ci)
		}

	case "enter":
		var sel []string
		for _, it := range l.Items() {
			if ci := it.(columnItem); ci.selected {
				sel = append(sel, ci.name)
			}
		}
		return sel
	}
	return nil
}

// applyFilter updates filteredRows based on current mode and input.
func (m *model) applyFilter() {
	p := m.textInput.Value()
	if m.regexMode {
		// regex
		re, err := regexp.Compile(p)
		if err != nil {
			// invalid—no change
			return
		}
		var out []Record
		for _, r := range m.allRows {
			for _, col := range m.filterCols {
				if re.MatchString(r[col]) {
					out = append(out, r)
					break
				}
			}
		}
		m.filteredRows = out
	} else {
		// substring
		lower := strings.ToLower(p)
		var out []Record
		for _, r := range m.allRows {
			if lower == "" {
				out = append(out, r)
				continue
			}
			for _, col := range m.filterCols {
				if strings.Contains(strings.ToLower(r[col]), lower) {
					out = append(out, r)
					break
				}
			}
		}
		m.filteredRows = out
	}
	m.buildTable(m.filteredRows)
}

// ─── Update & View ────────────────────────────────────────────────────────────

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	k := key.String()

	// global quit
	if isKey && (k == "ctrl+c" || k == "q") {
		return m, tea.Quit
	}

	switch m.state {

	// ─────────────── substring/regex filtering ───────────────
	case stateNavigation:
		if isKey {
			switch k {
			case "/":
				// search in non-regex mode
				m.regexMode = false
				m.state = stateFiltering
				m.applyFilter()
				return m, nil
			case "?":
				// search in non regex mode
				m.regexMode = true
				m.state = stateFiltering
				m.applyFilter()
				return m, nil
			case ",":
				if !m.regexMode {
					m.state = stateSelectVisibleColumns
					return m, nil
				}
			case ".":
				if !m.regexMode {
					m.state = stateSelectFilterColumns
					return m, nil
				}
			case "left", "h":
				m.table = m.table.ScrollLeft()
				return m, nil
			case "right", "l":
				m.table = m.table.ScrollRight()
				return m, nil
			case "up", "k", "down", "j":
				m.table, _ = m.table.Update(msg)
				return m, nil
			case "esc":
				m.table.Focused(true)
			}

		}

	case stateFiltering:
		m.textInput.Focus()

		if isKey {
			switch k {
			case "esc":
				m.table.Focused(true)
				m.textInput.Reset()
				m.applyFilter()
				m.state = stateNavigation
				return m, nil

			case "enter":
				m.table.Focused(true)
				m.applyFilter()
				m.state = stateNavigation
				return m, nil
			}
		}
		m.textInput, _ = m.textInput.Update(msg)
		m.applyFilter()
		return m, nil

	// ─────────────── pick visible columns ─────────────
	case stateSelectVisibleColumns:
		var cmd tea.Cmd
		m.listVisible, cmd = m.listVisible.Update(msg)
		if isKey {
			switch k {
			case "esc":
				m.table.Focused(true)
				m.textInput.Reset()
				m.applyFilter()
				m.state = stateNavigation
				return m, nil

			default:
				if sel := columnList(&m.listVisible, k); sel != nil {
					m.visibleCols = sel
					m.state = stateNavigation
					m.applyFilter()
					return m, nil
				}
			}
		}
		return m, cmd

	// ─────────────── pick filter columns ──────────────
	case stateSelectFilterColumns:
		var cmd tea.Cmd
		m.listFilter, cmd = m.listFilter.Update(msg)
		if isKey {

			switch k {
			case "esc":
				m.table.Focused(true)
				m.textInput.Reset()
				m.applyFilter()
				m.state = stateNavigation
				return m, nil
			default:
				if sel := columnList(&m.listFilter, k); sel != nil {
					m.filterCols = sel
					m.state = stateNavigation
					return m, nil
				}
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	mode := "substring"
	if m.regexMode {
		mode = "regex"
	}
	switch m.state {
	case stateFiltering:
		return fmt.Sprintf(
			"%s [%s]\n\n%s",
			m.textInput.View(),
			mode,
			m.table.View(),
		)
	case stateNavigation:
		return fmt.Sprintf(
			"%s [%s]\n%s\n%s",
			m.textInput.View(),
			mode,
			m.help.View(m.keys),
			m.table.View(),
		)
	case stateSelectVisibleColumns:
		return m.listVisible.View()
	case stateSelectFilterColumns:
		return m.listFilter.View()
	}
	return ""
}
