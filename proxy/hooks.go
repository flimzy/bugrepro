package errsql

import (
	"database/sql/driver"
	"io"
)

// Hooks is a collection of hooks which can be used to intercept and modify
// various events in the driver's operation.
type Hooks struct {
	// ErrorHook is called for any error, before it is returned to the caller.
	// The hook may return a modified error, or nil to indicate that the error
	// has been handled and should not be returned to the caller.
	ErrorHook func(*Event, error) error

	// BeforePrepare is called before a query is prepared with
	// [database/sql.DB.Prepare], [database/sql.DB.PrepareContext],
	// [database/sql.Conn.PrepareContext], [database/sql.Tx.Prepare], or
	// [database/sql.Tx.PrepareContext]. The hook may return a modified query
	// string. If the returned query string is empty, it is not modified. If an
	// error is returned, the query is aborted and the error is returned to the
	// caller, after being passed through ErrorHook, if defined.
	BeforePrepare func(*Event, string) (string, error)

	// BeforePreparedQuery is called before a prepared query is executed with
	// the [database/sql.Stmt.Query]. The hook may return a modified list of
	// arguments. If the returned list is nil, it is not modified. If an error
	// is returned, the query is aborted and the error is returned to the
	// caller, after being passed through ErrorHook, if defined.
	BeforePreparedQuery func(*Event, []driver.Value) ([]driver.Value, error)

	// BeforePreparedQueryContext is called before a prepared query is executed
	// with [database/sql.Stmt.QueryContext] or
	// [database/sql.Stmt.ExecContext]. The hook may return a modified list of
	// arguments. If the returned list is nil, it is not modified. If an error
	// is returned, the query is aborted and the error is returned to the
	// caller, after being passed through ErrorHook, if defined.
	BeforePreparedQueryContext func(*Event, []driver.NamedValue) ([]driver.NamedValue, error)

	// BeforeQuery is called before a query is executed with
	// [database/sql.DB.Query], [database/sql.DB.Exec],
	// [database/sql.Conn.Exec], [database/sql.Tx.Query], or
	// [database/sql.Tx.Exec]. The hook may return a modified query string or
	// arguments. If the returned query string is empty, it is not modified. If
	// the returned argument list is nil, the arguments are not modified. If an
	// error is returned, the query is aborted and the error is returned to the
	// caller, after being passed through ErrorHook, if defined.
	BeforeQuery func(*Event, string, []driver.Value) (string, []driver.Value, error)

	// BeforeQueryContext is called before a query is executed with
	// [database/sql.DB.QueryContext], [database/sql.DB.ExecContext]
	// [database/sql.Conn.QueryContext], [database/sql.Conn.ExecContext],
	// [database/sql.Tx.QueryContext], or [database/sql.Tx.ExecContext]. The
	// hook may return a modified query string or arguments. If the returned
	// query string is empty, it is not modified. If the returned argument list
	// is nil, the arguments are not modified. If an error is returned, the
	// query is aborted and the error is returned to the caller, after being
	// passed through ErrorHook, if defined.
	BeforeQueryContext func(*Event, string, []driver.NamedValue) (string, []driver.NamedValue, error)
}

// ErrorHandler is a function which is called for any error.
type ErrorHandler func(err error) error

func (d *Driver) beforeError(event *Event, err error) error {
	if d.hooks.ErrorHook == nil {
		return err
	}
	switch err {
	// Don't modify driver.Err* or io.EOF values
	case nil, driver.ErrSkip, driver.ErrBadConn, driver.ErrRemoveArgument, io.EOF:
		return err
	}
	return d.hooks.ErrorHook(event, err)
}

func (d *Driver) beforePrepare(event *Event, query string) (string, error) {
	if d.hooks.BeforePrepare == nil {
		return query, nil
	}
	q, err := d.hooks.BeforePrepare(event, query)
	if err != nil {
		return "", err
	}
	if q == "" {
		return query, nil
	}
	return q, nil
}

func (d *Driver) beforePreparedQuery(event *Event, args []driver.Value) ([]driver.Value, error) {
	if d.hooks.BeforePreparedQuery == nil {
		return args, nil
	}
	a, err := d.hooks.BeforePreparedQuery(event, args)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return args, nil
	}
	return a, nil
}

func (d *Driver) beforePreparedQueryContext(event *Event, args []driver.NamedValue) ([]driver.NamedValue, error) {
	if d.hooks.BeforePreparedQueryContext == nil {
		return args, nil
	}
	a, err := d.hooks.BeforePreparedQueryContext(event, args)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return args, nil
	}
	return a, nil
}

func (d *Driver) beforeQuery(event *Event, query string, args []driver.Value) (string, []driver.Value, error) {
	if d.hooks.BeforeQuery == nil {
		return query, args, nil
	}
	q, a, err := d.hooks.BeforeQuery(event, query, args)
	if err != nil {
		return "", nil, err
	}
	if q == "" {
		q = query
	}
	if a == nil {
		a = args
	}
	return q, a, nil
}

func (d *Driver) beforeQueryContext(event *Event, query string, args []driver.NamedValue) (string, []driver.NamedValue, error) {
	if d.hooks.BeforeQueryContext == nil {
		return query, args, nil
	}
	q, a, err := d.hooks.BeforeQueryContext(event, query, args)
	if err != nil {
		return "", nil, err
	}
	if q == "" {
		q = query
	}
	if a == nil {
		a = args
	}
	return q, a, nil
}
