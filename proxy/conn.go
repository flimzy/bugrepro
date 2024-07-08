package errsql

import (
	"database/sql/driver"
)

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
