package errsql_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"

	"gitlab.com/flimzy/errsql"
)

func Test_error_handlers(t *testing.T) {
	tests := []struct {
		name string
		eh   func(error) error
		fn   func(*testing.T, *sql.DB) error
		want string
	}{
		{
			name: "invalid syntax",
			eh:   errsql.AddStacktrace,
			fn: func(t *testing.T, db *sql.DB) error {
				_, err := db.Exec("THIS IS NOT VALID SQL")
				return err
			},
			want: `^near \"THIS\": syntax error\ngitlab.com/flimzy/errsql_test.Test_error_handlers`,
		},
		{
			name: "wrapped stack trace handler",
			eh: func(err error) error {
				return errsql.AddStacktrace(err)
			},
			fn: func(t *testing.T, db *sql.DB) error {
				_, err := db.Exec("THIS IS NOT VALID SQL")
				return err
			},
			want: `^near \"THIS\": syntax error\ngitlab.com/flimzy/errsql_test.Test_error_handlers.func3`, // This should ref func3, fn(), not func2, eh().
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := errsql.New(&sqlite3.SQLiteDriver{}, tt.eh)

			c, err := d.OpenConnector("file::memory:?mode=memory")
			if err != nil {
				t.Fatal(err)
			}
			db := sql.OpenDB(c)
			err = tt.fn(t, db)
			out := fmt.Sprintf("%+v", err)
			re := regexp.MustCompile(tt.want)
			if !re.MatchString(out) {
				t.Errorf("got %q, want %q", out, tt.want)
			}
		})
	}
}
