package middleware

type contextKey uint8

const (
	kContextValueRequestID contextKey = iota
	kContextLoggerStartTime
)
