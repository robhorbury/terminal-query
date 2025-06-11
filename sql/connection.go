package sql

import (
	"database/sql"
	"log/slog"

	dbsql "github.com/databricks/databricks-sql-go"
)

type Connection interface {
	Query(sqlString string) (*sql.Rows, error)
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
