package errsql

import (
	"database/sql/driver"
	"reflect"
)

func (c *connWrapper) newRows(r driver.Rows) driver.Rows {
	return &rowsWrapper{
		r: r,
		c: c,
	}
}

type rowsWrapper struct {
	r driver.Rows
	c *connWrapper
}

var (
	_ driver.Rows                           = (*rowsWrapper)(nil)
	_ driver.RowsNextResultSet              = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeScanType         = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeLength           = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeNullable         = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*rowsWrapper)(nil)
)

func (r *rowsWrapper) newEvent(method string) *Event {
	return newEvent(EntityRows, method, r.c.inTx)
}

func (r *rowsWrapper) Columns() []string {
	return r.r.Columns()
}

func (r *rowsWrapper) Close() error {
	return r.c.d.beforeError(r.newEvent(MethodClose), r.r.Close())
}

func (r *rowsWrapper) Next(dest []driver.Value) error {
	return r.c.d.beforeError(r.newEvent(MethodNext), r.r.Next(dest))
}

func (r *rowsWrapper) HasNextResultSet() bool {
	if s, ok := r.r.(driver.RowsNextResultSet); ok {
		return s.HasNextResultSet()
	}
	return false
}

func (r *rowsWrapper) NextResultSet() error {
	if s, ok := r.r.(driver.RowsNextResultSet); ok {
		return r.c.d.beforeError(r.newEvent(MethodNextResultSet), s.NextResultSet())
	}
	return nil
}

func (r *rowsWrapper) ColumnTypeScanType(index int) reflect.Type {
	if s, ok := r.r.(driver.RowsColumnTypeScanType); ok {
		return s.ColumnTypeScanType(index)
	}
	// fall back to default behavior of stdlib
	// See https://github.com/golang/go/blob/12e9b968bc8890b072b98facbd079ca337bd33a0/src/database/sql/sql.go#L3214
	return reflect.TypeOf(new(any)).Elem()
}

func (r *rowsWrapper) ColumnTypeDatabaseTypeName(index int) string {
	if s, ok := r.r.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return s.ColumnTypeDatabaseTypeName(index)
	}
	// Fall back to default behavior of stdlib
	// See https://github.com/golang/go/blob/12e9b968bc8890b072b98facbd079ca337bd33a0/src/database/sql/sql.go#L3216-L3218
	return ""
}

func (r *rowsWrapper) ColumnTypeLength(index int) (length int64, ok bool) {
	if s, ok := r.r.(driver.RowsColumnTypeLength); ok {
		return s.ColumnTypeLength(index)
	}
	// Fall back to default behavior of stdlib
	// See https://github.com/golang/go/blob/12e9b968bc8890b072b98facbd079ca337bd33a0/src/database/sql/sql.go#L3219-L3221
	return 0, false
}

func (r *rowsWrapper) ColumnTypeNullable(index int) (nullable, ok bool) {
	if s, ok := r.r.(driver.RowsColumnTypeNullable); ok {
		return s.ColumnTypeNullable(index)
	}
	// Fall back to default behavior of stdlib
	// See https://github.com/golang/go/blob/12e9b968bc8890b072b98facbd079ca337bd33a0/src/database/sql/sql.go#L3222-L3224
	return false, false
}

func (r *rowsWrapper) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	if s, ok := r.r.(driver.RowsColumnTypePrecisionScale); ok {
		return s.ColumnTypePrecisionScale(index)
	}
	// Fall back to default behavior of stdlib
	// See https://github.com/golang/go/blob/12e9b968bc8890b072b98facbd079ca337bd33a0/src/database/sql/sql.go#L3225-L3227
	return 0, 0, false
}
