package sql

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	dbsql "github.com/databricks/databricks-sql-go"
)

type Connection interface {
	Query(sqlString string) (*sql.Rows, error)
	RunQueryFromFile(filePath string) ([]map[string]string, error)
}

type DatabricksConnection struct {
	AccessToken    string
	ServerHostname string
	HttpPath       string
	Logger         *slog.Logger
}

func (c DatabricksConnection) Query(sqlString string) (*sql.Rows, error) {
	connector, err := dbsql.NewConnector(
		dbsql.WithAccessToken(c.AccessToken),
		dbsql.WithServerHostname(c.ServerHostname),
		dbsql.WithPort(443),
		dbsql.WithHTTPPath(c.HttpPath),
	)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	rows, err := db.Query(sqlString)
	return rows, err
}

func (c DatabricksConnection) RunQueryFromFile(filePath string) ([]map[string]string, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}
	rows, err := c.Query(string(data))
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Map of column names to value.
	var maps = []map[string]string{}

	for rows.Next() {
		vals := make([]any, len(cols))
		for i := range cols {
			var ii any
			vals[i] = &ii
		}

		err := rows.Scan(vals...)
		if err != nil {
			return nil, err
		}

		m := map[string]string{}
		for i, colName := range cols {
			raw := vals[i].(*any)
			m[colName] = fmt.Sprintf("%v", *raw)
		}

		maps = append(maps, m)
	}

	return maps, nil
}
