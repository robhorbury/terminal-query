package sql

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintRowsAsTableBasic writes rows in simple tab-delimited format
func PrintRowsAsTableBasic(writer io.Writer, rows []map[string]string) {
	if len(rows) == 0 {
		fmt.Fprintln(writer, "No rows to display.")
		return
	}
	cols := make([]string, 0, len(rows[0]))
	for c := range rows[0] {
		cols = append(cols, c)
	}
	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	// headers
	for _, c := range cols {
		fmt.Fprintf(w, "%s\t", c)
	}
	fmt.Fprintln(w)
	// rows
	for _, row := range rows {
		for _, c := range cols {
			fmt.Fprintf(w, "%v\t", row[c])
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}
