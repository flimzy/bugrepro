package errsql

import (
	"context"
	"database/sql/driver"
)

// newConn wraps c with.
func (d *Driver) newConn(c driver.Conn) driver.Conn {
	return &connWrapper{
		c: c,
		d: d,
	}
}

type connWrapper struct {
	d    *Driver
	c    driver.Conn
	inTx bool
}

var (
	_ driver.Conn               = (*connWrapper)(nil)
	_ driver.ConnBeginTx        = (*connWrapper)(nil)
	_ driver.ConnPrepareContext = (*connWrapper)(nil)
	_ driver.NamedValueChecker  = (*connWrapper)(nil)
	_ driver.Execer             = (*connWrapper)(nil)
	_ driver.ExecerContext      = (*connWrapper)(nil)
	// _ driver.Pinger             = (*connWrapper)(nil) // TODO
	_ driver.Queryer         = (*connWrapper)(nil)
	_ driver.QueryerContext  = (*connWrapper)(nil)
	_ driver.SessionResetter = (*connWrapper)(nil)
	_ driver.Validator       = (*connWrapper)(nil)
)

func (c *connWrapper) newEvent(method string) *Event {
	return newEvent(EntityConnection, method, c.inTx)
}

func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	event := c.newEvent(MethodPrepare)
	q, err := c.d.beforePrepare(event, query)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	stmt, err := c.c.Prepare(q)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newStmt(stmt), nil
}

func (c *connWrapper) Close() error {
	if err := c.c.Close(); err != nil {
		return c.d.beforeError(c.newEvent(MethodClose), err)
	}
	return nil
}

func (c *connWrapper) Begin() (driver.Tx, error) {
	tx, err := c.c.Begin()
	if err != nil {
		return nil, c.d.beforeError(c.newEvent(MethodBegin), err)
	}
	return c.newTx(tx), nil
}

func (c *connWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	connBeginTx, ok := c.c.(driver.ConnBeginTx)
	if !ok {
		return c.c.Begin()
	}
	tx, err := connBeginTx.BeginTx(ctx, opts)
	if err != nil {
		return nil, c.d.beforeError(c.newEvent(MethodBeginTx), err)
	}
	return c.newTx(tx), nil
}

func (c *connWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	event := c.newEvent(MethodPrepareContext)
	connPrepareContext, ok := c.c.(driver.ConnPrepareContext)
	if !ok {
		return c.c.Prepare(query)
	}
	q, err := c.d.beforePrepare(event, query)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	stmt, err := connPrepareContext.PrepareContext(ctx, q)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newStmt(stmt), nil
}

func (c *connWrapper) CheckNamedValue(nv *driver.NamedValue) error {
	namedValueChecker, ok := c.c.(driver.NamedValueChecker)
	if !ok {
		return driver.ErrSkip
	}
	return c.d.beforeError(c.newEvent(MethodCheckNamedValue), namedValueChecker.CheckNamedValue(nv))
}

func (c *connWrapper) Exec(query string, args []driver.Value) (driver.Result, error) {
	execer, ok := c.c.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}
	event := c.newEvent(MethodExec)
	q, a, err := c.d.beforeQuery(event, query, args)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	res, err := execer.Exec(q, a)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newResult(res), nil
}

func (c *connWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execerContext, ok := c.c.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	event := c.newEvent(MethodExecContext)
	q, a, err := c.d.beforeQueryContext(event, query, args)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	res, err := execerContext.ExecContext(ctx, q, a)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newResult(res), nil
}

func (c *connWrapper) Query(query string, args []driver.Value) (driver.Rows, error) {
	queryer, ok := c.c.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}
	event := c.newEvent(MethodQuery)
	q, a, err := c.d.beforeQuery(event, query, args)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	rows, err := queryer.Query(q, a)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newRows(rows), nil
}

func (c *connWrapper) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryerContext, ok := c.c.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	event := c.newEvent(MethodQueryContext)
	q, a, err := c.d.beforeQueryContext(event, query, args)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	rows, err := queryerContext.QueryContext(ctx, q, a)
	if err != nil {
		return nil, c.d.beforeError(event, err)
	}
	return c.newRows(rows), nil
}

func (c *connWrapper) ResetSession(ctx context.Context) error {
	sessionResetter, ok := c.c.(driver.SessionResetter)
	if !ok {
		return nil
	}
	return c.d.beforeError(c.newEvent(MethodResetSession), sessionResetter.ResetSession(ctx))
}

func (c *connWrapper) IsValid() bool {
	validator, ok := c.c.(driver.Validator)
	if !ok {
		return true
	}
	return validator.IsValid()
}
