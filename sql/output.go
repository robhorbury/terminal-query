package sql

import (
	"fmt"
	"io"
	"text/tabwriter"

	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func PrintRowsAsTableTea(input_rows []map[string]string) {

	// Extract and fix column order from first row
	columns := []table.Column{}
	column_names := []string{}

	for col := range input_rows[0] {
		columns = append(columns, table.Column{Title: col, Width: 10})
		column_names = append(column_names, col)
	}

	// Print rows
	rows := []table.Row{}
	for _, input_row := range input_rows {
		row := table.Row{}
		for _, col := range column_names {
			row = append(row, input_row[col])
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
