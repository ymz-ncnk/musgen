package musgen

import "errors"

var ErrUnsupportedEncoding = errors.New("unsupported encoding")
var ErrUnsupportedElemEncoding = errors.New("unsupported elem encoding")
var ErrUnsupportedKeyEncoding = errors.New("unsupported key encoding")
