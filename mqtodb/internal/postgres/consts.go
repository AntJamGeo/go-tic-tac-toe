package postgres

const (
	numDbConnectRetries = 5

	// Message Types
	register   = "register"
	deregister = "deregister"
	update     = "update"
)
