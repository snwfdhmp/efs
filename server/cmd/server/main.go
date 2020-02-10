package main

import (
	"flag"

	"github.com/snwfdhmp/efs/server/services"
	"go.uber.org/zap"
)

var (
	log, _     = zap.NewDevelopment()
	serverPort = flag.String("port", "8081", "http port")
	path       = flag.String("fs-path", "example/efs/user-data", "path for fs")
)

func init() {
	flag.Parse()
}

func main() {
	ctx, err := services.NewCtx(log, *path)
	if err != nil {
		log.Fatal("could not start server", zap.Error(err))
		return
	}

	log.Info("Starting API. Listening on " + *serverPort + ".")
	ctx.API().Listen("0.0.0.0:" + *serverPort)
}
