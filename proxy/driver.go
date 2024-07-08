package errsql

import (
	"database/sql/driver"
)

// Driver is the core type implementing the [database/sql/driver.Driver] interface.
type Driver struct {
	d     driver.Driver
	hooks *Hooks
}

var (
	_ driver.Driver        = (*Driver)(nil)
	_ driver.DriverContext = (*Driver)(nil)
)

func (*Driver) newEvent(method string) *Event {
	return newEvent(EntityDriver, method, false)
}

// New returns a new Driver instance, which wraps driver and calls eh for
// any errors. The default error handler just passes the error through,
// unaltered.
func New(driver driver.Driver, eh ErrorHandler) *Driver {
	return NewWithHooks(driver, &Hooks{
		ErrorHook: func(_ *Event, err error) error {
			return err
		},
	})
}

// NewWithHooks returns a new Driver instance, which wraps driver and calls
// any defined hooks for their respective events. See [Hooks] for an explanation
// of the available hooks.
func NewWithHooks(driver driver.Driver, hooks *Hooks) *Driver {
	return &Driver{
		d:     driver,
		hooks: hooks,
	}
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.d.Open(name)
	if err != nil {
		return nil, d.beforeError(d.newEvent(MethodOpen), err)
	}
	return d.newConn(conn), nil
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	driverContext, ok := d.d.(driver.DriverContext)
	if !ok {
		return &connectorShim{
			d:    d,
			name: name,
		}, nil
	}
	connector, err := driverContext.OpenConnector(name)
	if err != nil {
		return nil, d.beforeError(d.newEvent(MethodOpenConnector), err)
	}
	return d.newConnector(connector), nil
}
