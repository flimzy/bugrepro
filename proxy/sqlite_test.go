package errsql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
	"testing"
)

func newSQLite(t *testing.T, w io.Writer) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	var hooks *Hooks
	if w != nil {
		logger := log.New(w, "", 0)
		hooks = &Hooks{
			BeforePrepare: func(e *Event, query string) (string, error) {
				logger.Printf("[%s.%s]: %s", e.Entity, e.Method, query)
				return query, nil
			},
			BeforeQuery: func(e *Event, query string, args []driver.Value) (string, []driver.Value, error) {
				logger.Printf("[%s.%s]: %s %v", e.Entity, e.Method, query, args)
				return query, args, nil
			},
			BeforeQueryContext: func(e *Event, query string, args []driver.NamedValue) (string, []driver.NamedValue, error) {
				logger.Printf("[%s.%s]: %s %v", e.Entity, e.Method, query, args)
				return query, args, nil
			},
		}
	}
	drv := NewWithHooks(db.Driver(), hooks)
	connector, err := drv.OpenConnector(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return sql.OpenDB(connector)
}

func TestSQLite(t *testing.T) {
	t.Run("DB.Exec", func(t *testing.T) {
		log := &bytes.Buffer{}
		db := newSQLite(t, log)
		defer db.Close()
		_, err := db.Exec("SELECT 1")
		if err != nil {
			t.Fatal(err)
		}

		if got, want := log.String(), "[connection.ExecContext]: SELECT 1 []\n"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
