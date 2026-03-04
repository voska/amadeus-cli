package errfmt

import "fmt"

const (
	ExitOK        = 0
	ExitError     = 1
	ExitUsage     = 2
	ExitEmpty     = 3
	ExitAuth      = 4
	ExitNotFound  = 5
	ExitForbidden = 6
	ExitRateLimit = 7
	ExitRetryable = 8
	ExitConfig    = 10
)

type Error struct {
	Code    int
	Message string
	Detail  string
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func Wrap(code int, msg string, err error) *Error {
	return &Error{Code: code, Message: msg, Detail: err.Error()}
}

func Auth(msg string) *Error      { return New(ExitAuth, msg) }
func NotFound(msg string) *Error  { return New(ExitNotFound, msg) }
func Usage(msg string) *Error     { return New(ExitUsage, msg) }
func Empty() *Error               { return New(ExitEmpty, "no results") }
func Config(msg string) *Error    { return New(ExitConfig, msg) }
func RateLimit() *Error           { return New(ExitRateLimit, "rate limited") }
func Forbidden(msg string) *Error { return New(ExitForbidden, msg) }
