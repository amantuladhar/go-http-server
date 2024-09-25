package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/pkg/cliargs"
	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
	"github.com/codecrafters-io/http-server-starter-go/pkg/zhttp"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	argsWithoutProg := os.Args[1:]
	util.LogDebug(fmt.Sprintf("%v", argsWithoutProg))

	config := zhttp.NewHttpServerConfig()
	config.HandleFunc("/", root)
	config.HandleFunc("/echo/{slug}", echo)
	config.HandleFunc("/user-agent", userAgent)
	config.HandleFunc("/files/{filename}", fileName)

	err := zhttp.ListenAndServe("0.0.0.0:4221", config)
	util.ExitOnErr(err, "unable to start server")
}
func fileName(r *zhttp.Request) *zhttp.Response {
	content, err := os.ReadFile(fmt.Sprintf("/%s/%s", cliargs.GetArg("--directory"), r.PathParam["filename"]))
	if err != nil {
		return zhttp.NewResponse().StatusCode(404)
	}
	return zhttp.NewResponse().File(content).StatusCode(200)
}

func userAgent(r *zhttp.Request) *zhttp.Response {
	return zhttp.NewResponse().Text([]byte(r.Headers["User-Agent"])).StatusCode(200)
}

func echo(r *zhttp.Request) *zhttp.Response {
	return zhttp.NewResponse().Text([]byte(r.PathParam["slug"])).StatusCode(200)
}

func root(r *zhttp.Request) *zhttp.Response {
	return zhttp.NewResponse().StatusCode(200)
}
