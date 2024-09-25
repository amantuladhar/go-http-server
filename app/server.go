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
	config.HandleFunc("GET /", root)
	config.HandleFunc("GET /echo/{slug}", echo)
	config.HandleFunc("GET /user-agent", userAgent)
	config.HandleFunc("GET /files/{filename}", fileName)
	config.HandleFunc("POST /files/{filename}", postFileName)

	err := zhttp.ListenAndServe("0.0.0.0:4221", config)
	util.ExitOnErr(err, "unable to start server")
}

func postFileName(r *zhttp.Request) *zhttp.Response {
	f, err := os.Create(fmt.Sprintf("/%s/%s", cliargs.GetArg("--directory"), r.PathParam["filename"]))

	if err != nil {
		return zhttp.NewResponse().Text([]byte("unable to create file")).StatusCode(500)
	}
	_, err = f.Write(r.Body)
	if err != nil {
		return zhttp.NewResponse().Text([]byte("unable to write")).StatusCode(500)
	}
	return zhttp.NewResponse().StatusCode(201)
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
