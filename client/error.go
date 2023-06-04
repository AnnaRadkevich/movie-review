package client

import (
	"fmt"
	"net/http"
)

var _ error = (*Error)(nil)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.Code), e.Message)
}
