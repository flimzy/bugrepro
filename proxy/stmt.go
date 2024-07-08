package errsql

import (
	"context"
	"database/sql/driver"
	"errors"
)

// newStmt wraps s.
func (c *connWrapper) newStmt(s driver.Stmt) driver.Stmt {
	return &stmtWrapper{
		s: s,
		c: c,
	}
}

type stmtWrapper struct {
	c *connWrapper
	s driver.Stmt
}

var (
	_ driver.Stmt              = (*stmtWrapper)(nil)
	_ driver.NamedValueChecker = (*stmtWrapper)(nil)
	_ driver.ColumnConverter   = (*stmtWrapper)(nil)
	_ driver.StmtExecContext   = (*stmtWrapper)(nil)
	_ driver.StmtQueryContext  = (*stmtWrapper)(nil)
)

func (s *stmtWrapper) newEvent(method string) *Event {
	return newEvent(EntityStatement, method, s.c.inTx)
}

func (s *stmtWrapper) Close() error {
	return s.c.d.beforeError(s.newEvent(MethodClose), s.s.Close())
}

func (s *stmtWrapper) NumInput() int {
	return s.s.NumInput()
}

func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	event := s.newEvent(MethodExec)
	a, err := s.c.d.beforePreparedQuery(event, args)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	res, err := s.s.Exec(a)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	return s.c.newResult(res), nil
}

func (s *stmtWrapper) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	event := s.newEvent(MethodExecContext)
	execContext, ok := s.s.(driver.StmtExecContext)
	if !ok {
		dargs, err := namedValueToValue(args)
		if err != nil {
			return nil, s.c.d.beforeError(event, err)
		}
		return s.Exec(dargs)
	}
	a, err := s.c.d.beforePreparedQueryContext(event, args)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	res, err := execContext.ExecContext(ctx, a)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	return s.c.newResult(res), nil
}

func (s *stmtWrapper) Query(args []driver.Value) (driver.Rows, error) {
	event := s.newEvent(MethodQuery)
	a, err := s.c.d.beforePreparedQuery(event, args)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	rows, err := s.s.Query(a)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	return s.c.newRows(rows), nil
}

func (s *stmtWrapper) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	event := s.newEvent(MethodQueryContext)
	queryContext, ok := s.s.(driver.StmtQueryContext)
	if !ok {
		dargs, err := namedValueToValue(args)
		if err != nil {
			return nil, s.c.d.beforeError(event, err)
		}
		return s.Query(dargs)
	}
	a, err := s.c.d.beforePreparedQueryContext(event, args)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	rows, err := queryContext.QueryContext(ctx, a)
	if err != nil {
		return nil, s.c.d.beforeError(event, err)
	}
	return s.c.newRows(rows), nil
}

func (s *stmtWrapper) CheckNamedValue(nv *driver.NamedValue) error {
	namedValueChecker, ok := s.s.(driver.NamedValueChecker)
	if !ok {
		return driver.ErrSkip
	}
	return s.c.d.beforeError(s.newEvent(MethodCheckNamedValue), namedValueChecker.CheckNamedValue(nv))
}

func (s *stmtWrapper) ColumnConverter(idx int) driver.ValueConverter {
	columnConverter, ok := s.s.(driver.ColumnConverter)
	if !ok {
		return driver.DefaultParameterConverter
	}
	return s.newValuesConverter(columnConverter.ColumnConverter(idx))
}

// copied from stdlib
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
