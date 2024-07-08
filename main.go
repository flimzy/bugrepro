package main

import (
	"context"
	"database/sql"
	"database/sql/driver"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	drv := db.Driver()
	edrv := New(drv)
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

// Driver is the core type implementing the [database/sql/driver.Driver] interface.
type Driver struct {
	driver.Driver
}

var _ driver.Driver = (*Driver)(nil)

// New returns a new Driver instance, which wraps driver and calls eh for
// any errors. The default error handler just passes the error through,
// unaltered.
func New(driver driver.Driver) *Driver {
	return &Driver{Driver: driver}
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return driver.Conn(&connWrapper{Conn: conn}), nil
}

type connWrapper struct {
	driver.Conn
}

var _ driver.Conn = (*connWrapper)(nil)

func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return driver.Stmt(&stmtWrapper{Stmt: stmt}), nil
}

type stmtWrapper struct {
	driver.Stmt
}

var (
	_ driver.Stmt            = (*stmtWrapper)(nil)
	_ driver.ColumnConverter = (*stmtWrapper)(nil)
)

func (s *stmtWrapper) ColumnConverter(idx int) driver.ValueConverter {
	columnConverter, ok := s.Stmt.(driver.ColumnConverter)
	if !ok {
		return driver.DefaultParameterConverter
	}
	return driver.ValueConverter(&valueConverterWrapper{v: columnConverter.ColumnConverter(idx)})
}

type valueConverterWrapper struct {
	v driver.ValueConverter
}

var _ driver.ValueConverter = valueConverterWrapper{}

func (vc valueConverterWrapper) ConvertValue(v any) (driver.Value, error) {
	return vc.v.ConvertValue(v)
}
