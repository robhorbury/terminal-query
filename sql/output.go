package sql

import (
	"fmt"
	"io"
	"text/tabwriter"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

func PrintRowsAsTableBasic(writer io.Writer, rows []map[string]string) {
	if len(rows) == 0 {
		fmt.Fprintln(writer, "No rows to display.")
		return
	}

	// Extract and fix column order from first row
	columns := make([]string, 0, len(rows[0]))
	for col := range rows[0] {
		columns = append(columns, col)
	}

	// Use tabwriter
	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)

	// Print headers
	for _, col := range columns {
		fmt.Fprintf(w, "%s\t", col)
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range rows {
		for _, col := range columns {
			fmt.Fprintf(w, "%v\t", row[col])
		}
		fmt.Fprintln(w)
	}

	// Flush the writer
	w.Flush()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			// Toggle focus
			m.table = m.table.Focused(!m.table.GetFocused())

		case "left":
			m.table = m.table.ScrollLeft()
			return m, nil

		case "h":
			m.table = m.table.ScrollLeft()
			return m, nil

		case "right":
			m.table = m.table.ScrollRight()
			return m, nil

		case "l":
			m.table = m.table.ScrollRight()
			return m, nil

		case "enter":
			sel := m.table.SelectedRows()
			if len(sel) > 0 {
				fmt.Printf("Selected row data: %#v\n", sel[0].Data)
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.table.View() + "\n"
}

func PrintRowsAsTableTea(input []map[string]string, input_columns []string) {
	if len(input) == 0 {
		fmt.Println("No data to display.")
		return
	}

	// Detect terminal width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 120 // fallback width
	}

	// Build columns and keys
	colKeys := []string{}
	columns := []table.Column{}
	for _, k := range input_columns {
		colKeys = append(colKeys, k)
		columns = append(columns, table.NewColumn(k, k, 25)) // wider default width
	}

	// Build rows
	rows := make([]table.Row, len(input))
	for i, rec := range input {
		rd := table.RowData{}
		for _, k := range colKeys {
			rd[k] = rec[k]
		}
		rows[i] = table.NewRow(rd)
	}
	// Create table
	tbl := table.New(columns).
		WithRows(rows).
		WithMaxTotalWidth(width).
		WithPageSize(25). // ← limit rows per page
		WithHorizontalFreezeColumnCount(1).
		WithMinimumHeight(10).
		Focused(true).
		WithBaseStyle(baseStyle).
		HeaderStyle(
			lipgloss.NewStyle().
				Bold(true).
				Underline(true).
				Foreground(lipgloss.Color("250")),
		).
		HighlightStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")),
		)

	m := model{table: tbl}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
