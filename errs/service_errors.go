package errs

import "errors"

var EventExistError error
var EventNotExistError error

func init() {
	EventExistError = errors.New("event has already existed")
	EventNotExistError = errors.New("event does not exist")
}
