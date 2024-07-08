package errsql

import "database/sql/driver"

func (c *connWrapper) newResult(r driver.Result) driver.Result {
	return &resultWrapper{
		r: r,
		c: c,
	}
}

type resultWrapper struct {
	r driver.Result
	c *connWrapper
}

var _ driver.Result = (*resultWrapper)(nil)

func (r *resultWrapper) newEvent(method string) *Event {
	return newEvent(EntityResult, method, r.c.inTx)
}

func (r *resultWrapper) LastInsertId() (int64, error) {
	id, err := r.r.LastInsertId()
	return id, r.c.d.beforeError(r.newEvent(MethodLastInsertId), err)
}

func (r *resultWrapper) RowsAffected() (int64, error) {
	ra, err := r.r.RowsAffected()
	return ra, r.c.d.beforeError(r.newEvent(MethodRowsAffected), err)
}
