package domain

type ContextKey string

const (
	LoggerKey    ContextKey = "logger"
	RequestIDKey ContextKey = "request_id"
)
