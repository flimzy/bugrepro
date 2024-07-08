package main

import (
	"context"
	"database/sql"

	"gitlab.com/flimzy/errsql"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	drv := db.Driver()
	connector, err := errsql.NewWithHooks(drv, &errsql.Hooks{}).OpenConnector(":memory:")
	if err != nil {
		panic(err)
	}
	db = sql.OpenDB(connector)

	stmt, err := db.PrepareContext(context.Background(), "SELECT $1")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(1)
	if err != nil {
		panic(err)
	}
}
