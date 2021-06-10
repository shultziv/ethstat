package etherscan

import (
	"fmt"
)

type errInvalidResponseFromApi struct {
	Url      string
	Response []byte
}

func (e *errInvalidResponseFromApi) Error() string {
	return fmt.Sprintf("Error ivalid response from api, request to URL: %s. Response: %v", e.Url, e.Response)
}

type errAccessToApi struct {
	Url        string
	ErrMessage string
}

func (e *errAccessToApi) Error() string {
	return fmt.Sprintf("Error access to api, request to URL: %s. Error message: %s", e.Url, e.ErrMessage)
}

type errNotFoundRT struct{}

func (e *errNotFoundRT) Error() string {
	return fmt.Sprintf("Not found round trippers")
}
