package musgen

import (
	"errors"
)

// ErrUnsupportedEncoding happens when MusGen tries to generate code for a
// field description with unsupported encoding.
var ErrUnsupportedEncoding = errors.New("unsupported encoding")

// ErrUnsupportedElemEncoding happens when MusGen tries to generate code for a 
// field description with unsupported elem encoding.
var ErrUnsupportedElemEncoding = errors.New("unsupported elem encoding")

// ErrUnsupportedKeyEncoding happens when MusGen tries to generate code for a 
// field description with unsupported key encoding.
var ErrUnsupportedKeyEncoding = errors.New("unsupported key encoding")

// ErrUnsupportedLand happens when MusGen tries to generate code for unsupported
// language.
var ErrUnsupportedLang = errors.New("unsupported lang")

