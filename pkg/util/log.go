package util

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func LogInfo(msg string, args ...interface{}) {
	_, f, l, ok := runtime.Caller(1)
	file := "unknown"
	line := 0
	if ok {
		file = f
		line = l
	}
	slog.Info(msg, append(args, "file", fmt.Sprintf(" %s:%d", file, line))...)
}

func LogDebug(msg string, args ...interface{}) {
	_, f, l, ok := runtime.Caller(1)
	file := "unknown"
	line := 0
	if ok {
		file = f
		line = l
	}
	slog.Debug(msg, append(args, "file", fmt.Sprintf(" %s:%d", file, line))...)
}

func LogErr(msg string, args ...interface{}) {
	_, f, l, ok := runtime.Caller(1)
	file := "unknown"
	line := 0
	if ok {
		file = f
		line = l
	}
	slog.Error(msg, append(args, "file", fmt.Sprintf(" %s:%d", file, line))...)
}

func LogOnErr(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	_, f, l, ok := runtime.Caller(1)
	file, line := "unknown", 0
	if ok {
		file = f
		line = l
	}
	slog.Error(msg, append(args, "err", err, "file", fmt.Sprintf(" %s:%d", file, line))...)
}

func ExitOnErr(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	_, f, l, ok := runtime.Caller(1)
	file, line := "unknown", 0
	if ok {
		file = f
		line = l
	}
	slog.Error(msg, append(args, "err", err, "file", fmt.Sprintf(" %s:%d", file, line))...)
	os.Exit(1)
}
