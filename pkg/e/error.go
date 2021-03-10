package e

import "errors"

var ErrorMsg = map[string]string{

	"RECORD_NOT_FOUND": 				 "record not found",
}

func GetError(s string) error {
	return errors.New(ErrorMsg[s])
}