package util

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func ExitOnErr(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	_, f, l, ok := runtime.Caller(1)
	if ok {
		slog.Error(msg, append(args, "err", err, "file", fmt.Sprintf(" %s:%d", f, l))...)
	} else {
		slog.Error(msg, append(args, "err", err)...)
	}
	os.Exit(1)
}
