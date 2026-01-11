package errors

import (
	"errors"
	"fmt"
)

var ErrorUserNotFound = errors.New(fmt.Sprint("user id not found"))
var ErrDataNotFound = errors.New("data not found")
