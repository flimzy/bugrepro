package main

import (
	"context"
	"database/sql"

	errsql "github.com/flimzy/bugrepro/proxy"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	drv := db.Driver()
	edrv := errsql.New(drv)
	sql.Register("proxy", edrv)
	db, err = sql.Open("proxy", ":memory:")
	if err != nil {
		panic(err)
	}

	stmt, err := db.PrepareContext(context.Background(), "SELECT $1")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(1)
	if err != nil {
		panic(err)
	}
}
