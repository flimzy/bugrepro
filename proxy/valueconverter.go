package errsql

import "database/sql/driver"

type valueConverterWrapper struct {
	v driver.ValueConverter
}

var _ driver.ValueConverter = valueConverterWrapper{}

func (vc valueConverterWrapper) ConvertValue(v any) (driver.Value, error) {
	return vc.v.ConvertValue(v)
}
