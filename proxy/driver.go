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
