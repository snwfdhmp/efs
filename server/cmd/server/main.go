package main

import (
	"github.com/snwfdhmp/efs/server/services"
	"go.uber.org/zap"
)

var (
	log, _     = zap.NewDevelopment()
	serverPort = "8081"
)

func main() {
	path := "example/efs/user-data"
	ctx, err := services.NewCtx(log, path)
	if err != nil {
		log.Fatal("could not start server", zap.Error(err))
		return
	}

	log.Info("Starting API. Listening on " + serverPort + ".")
	ctx.API().Listen("0.0.0.0:" + serverPort)
}
