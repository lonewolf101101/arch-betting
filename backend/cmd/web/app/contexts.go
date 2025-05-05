package app

type contextKey string

const (
	ContextKeyIsAuth       = contextKey("isAuthenticated")
	ContextKeyAuthCustomer = contextKey("authenticatedCustomer")
)
