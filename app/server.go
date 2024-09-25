package main

import (
	"log/slog"

	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
	"github.com/codecrafters-io/http-server-starter-go/pkg/zhttp"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	config := zhttp.NewHttpServerConfig()
	config.HandleFunc("/", root)
	config.HandleFunc("/echo/{slug}", echo)
	config.HandleFunc("/user-agent", userAgent)

	err := zhttp.ListenAndServe("0.0.0.0:4221", config)
	util.ExitOnErr(err, "unable to start server")
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
