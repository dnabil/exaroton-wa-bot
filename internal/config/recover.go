package config

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

// Is used with defer.
// For the log implementation, refer to RecoverHelper.
func Recover(ctx context.Context, args *Args) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}

		ErrLog(ctx, args, err, debug.Stack())

		os.Exit(1)
	}
}
