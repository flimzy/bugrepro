package errsql

import (
	"context"
	"database/sql/driver"
)

// WrapConnector wraps a [database/sql/driver.Connector] so that any errors
// returned by its [driver.Connector.Connect] method are passed through the
// error handler. Use this in place of [New] if you already have a
// [database/sql/driver.Connector] instance that you wish to wrap.
func WrapConnector(c driver.Connector, eh ErrorHandler) driver.Connector {
	d := New(c.Driver(), eh)
	return &connectorWrapper{
		c: c,
		d: d,
	}
}

func (d *Driver) newConnector(c driver.Connector) driver.Connector {
	return &connectorWrapper{
		c: c,
		d: d,
	}
}

type connectorWrapper struct {
	c driver.Connector
	d *Driver
}

var _ driver.Connector = (*connectorWrapper)(nil)

func (*connectorWrapper) newEvent(method string) *Event {
	return newEvent(EntityConnector, method, false)
}

func (c *connectorWrapper) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.c.Connect(ctx)
	if err != nil {
		return nil, c.d.beforeError(c.newEvent(MethodConnect), err)
	}
	return c.d.newConn(conn), nil
}

func (c *connectorWrapper) Driver() driver.Driver {
	return c.d
}

// connectorShim implements the [database/sql/driver.Connector] interface
// for a driver that doesn't support it. It ignores the [context.Context]
// argument to its Connect method.
type connectorShim struct {
	d    *Driver
	name string
}

func (c *connectorShim) Connect(context.Context) (driver.Conn, error) {
	return c.d.Open(c.name)
}

func (c *connectorShim) Driver() driver.Driver {
	return c.d
}
