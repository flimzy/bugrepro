package errsql

import "database/sql/driver"

func (c *connWrapper) newTx(t driver.Tx) driver.Tx {
	c.inTx = true
	return &txWrapper{
		t: t,
		c: c,
	}
}

type txWrapper struct {
	t driver.Tx
	c *connWrapper
}

var _ driver.Tx = (*txWrapper)(nil)

func (t *txWrapper) newEvent(method string) *Event {
	return newEvent(EntityTransaction, method, t.c.inTx)
}

func (t *txWrapper) Commit() error {
	err := t.t.Commit()
	if err != nil {
		return t.c.d.beforeError(t.newEvent(MethodCommit), err)
	}
	t.c.inTx = false
	return nil
}

func (t *txWrapper) Rollback() error {
	err := t.t.Rollback()
	if err != nil {
		return t.c.d.beforeError(t.newEvent(MethodRollback), err)
	}
	t.c.inTx = false
	return nil
}
