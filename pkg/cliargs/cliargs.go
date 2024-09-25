package cliargs

import (
	"os"
	"sync"
)

var once sync.Once
var args map[string]string

func GetArg(key string) string {
	once.Do(func() {
		args = make(map[string]string)
		argsWithoutProg := os.Args[1:]
		for i := 0; i < len(argsWithoutProg); i += 2 {
			args[argsWithoutProg[i]] = argsWithoutProg[i+1]
		}
	})
	return args[key]
}
