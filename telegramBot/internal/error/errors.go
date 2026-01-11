package error

import (
	"errors"
	"fmt"
)

var ErrorUserNotFound = errors.New(fmt.Sprint("user id not found"))
