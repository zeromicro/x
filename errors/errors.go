package errors

import "fmt"

// CodeMsg is a struct that contains a code and a message.
// It implements the error interface.
type CodeMsg struct {
	Code int
	Msg  string
}

func (c *CodeMsg) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", c.Code, c.Msg)
}

// New creates a new CodeMsg.
func New(code int, msg string) error {
	return &CodeMsg{Code: code, Msg: msg}
}
