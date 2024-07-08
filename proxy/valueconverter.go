package errsql

import "database/sql/driver"

func (s *stmtWrapper) newValuesConverter(v driver.ValueConverter) driver.ValueConverter {
	return &valueConverterWrapper{
		v: v,
		c: s.c,
	}
}

type valueConverterWrapper struct {
	v driver.ValueConverter
	c *connWrapper
}

var _ driver.ValueConverter = valueConverterWrapper{}

func (vc *valueConverterWrapper) newEvent(method string) *Event {
	return newEvent(EntityValueConverter, method, vc.c.inTx)
}

func (vc valueConverterWrapper) ConvertValue(v any) (driver.Value, error) {
	v, err := vc.v.ConvertValue(v)
	return v, vc.c.d.beforeError(vc.newEvent(MethodConvertValue), err)
}
