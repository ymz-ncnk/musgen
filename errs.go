package musgen

import (
	"errors"
	"fmt"
)

var ErrUnsupportedEncoding = errors.New("unsupported encoding")
var ErrUnsupportedElemEncoding = errors.New("unsupported elem encoding")
var ErrUnsupportedKeyEncoding = errors.New("unsupported key encoding")

type UnsupportedLangError struct {
	lang Lang
}

func (err UnsupportedLangError) Lang() Lang {
	return err.lang
}

func (err UnsupportedLangError) Error() string {
	return fmt.Sprintf("unsupported lang %s", err.lang.String())
}
