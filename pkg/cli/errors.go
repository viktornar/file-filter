package cli

import "errors"

var (
	ErrMissingArguments = errors.New("cli: missing arguments")
)