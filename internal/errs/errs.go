package errs

import "errors"

// config errors
var (
	ErrCfgYmlPathNotSet = errors.New(".yml path is not set")
)
