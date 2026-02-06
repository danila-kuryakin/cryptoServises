package error

import (
	"errors"
	"fmt"
)

var ErrorUserNotFound = errors.New(fmt.Sprint("user id not found"))
var ErrorNotFound = errors.New(fmt.Sprint("record not found"))
