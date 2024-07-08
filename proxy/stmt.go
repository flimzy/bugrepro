package errsql

import (
	"database/sql/driver"
)

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
