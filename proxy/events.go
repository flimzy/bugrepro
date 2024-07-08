package errsql

const (
	// Connection refers to a [database/sql/driver.Conn] instance.
	EntityConnection = "connection"
	// Connector refers to a [database/sql/driver.Connector] instance.
	EntityConnector = "connector"
	// ConnectorShim refers to a [database/sql/driver.Connector] shim instance,
	// which is used to wrap a driver that doesn't support the Connector
	// interface.
	EntityConnectorShim = "connectorShim"
	// Result refers to a [database/sql/driver.Driver] instance.
	EntityDriver = "driver"
	// Result refers to a [database/sql/driver.Result] instance.
	EntityResult = "result"
	// Rows refers to a [database/sql/driver.Rows] instance.
	EntityRows = "rows"
	// Statement refers to a [database/sql/driver.Stmt] instance.
	EntityStatement = "statement"
	// Transaction refers to a [database/sql/driver.Tx] instance.
	EntityTransaction = "transaction"
	// ValueConverter refers to a [database/sql/driver.ValueConverter] instance.
	EntityValueConverter = "valueConverter"

	MethodBegin           = "Begin"
	MethodBeginTx         = "BeginTx"
	MethodCheckNamedValue = "CheckNamedValue"
	MethodClose           = "Close"
	MethodCommit          = "Commit"
	MethodConnect         = "Connect"
	MethodConvertValue    = "ConvertValue"
	MethodExec            = "Exec"
	MethodExecContext     = "ExecContext"
	MethodLastInsertId    = "LastInsertId"
	MethodNext            = "Next"
	MethodNextResultSet   = "NextResultSet"
	MethodOpen            = "Open"
	MethodOpenConnector   = "OpenConnector"
	MethodPrepare         = "Prepare"
	MethodPrepareContext  = "PrepareContext"
	MethodQuery           = "Query"
	MethodQueryContext    = "QueryContext"
	MethodResetSession    = "ResetSession"
	MethodRollback        = "Rollback"
	MethodRowsAffected    = "RowsAffected"
)

type Event struct {
	// Entity is the entity which the event is related to, e.g. "connection",
	// "statement", "result", etc.
	Entity string
	// Method is the method which the event is related to, e.g. "Prepare",
	// "Exec", "Query", etc.
	Method string

	// InTransaction is true if the current connection was in a transaction
	// at the time of the event.
	InTransaction bool
}

func newEvent(entity, method string, inTx bool) *Event {
	return &Event{
		Entity:        entity,
		Method:        method,
		InTransaction: inTx,
	}
}
