package errsql

import (
	"database/sql/driver"
)

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
