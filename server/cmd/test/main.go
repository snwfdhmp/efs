package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/snwfdhmp/efs/server/services"
	"go.uber.org/zap"
)

var (
	log, _ = zap.NewDevelopment()
)

func main() {
	ctx, err := TestCreateCtx()
	if err != nil {
		log.Error("fatal:", zap.Error(err))
		return
	}

	path, err := TestPostFileTmp(ctx)
	if err != nil {
		log.Error("fatal:", zap.Error(err))
		return
	}

	content, err := TestGetFileTmp(ctx, path)
	if err != nil {
		log.Error("fatal:", zap.Error(err))
		return
	}

	log.Info("content: " + string(content))

	log.Info("TESTS SUCCESSFUL")
}

func TestCreateCtx() (services.Ctx, error) {
	path := "./example/efs/user-data"
	log.Info("Tests> Testing services.NewCtx...")
	return services.NewCtx(log, path)
}

func TestPostFileTmp(ctx services.Ctx) (path string, err error) {
	path = fmt.Sprintf("/tmp/testfile-%d-%d", time.Now().Unix(), rand.Intn(100000))
	err = ctx.PostFile(path, fileContent)
	return
}

func TestGetFileTmp(ctx services.Ctx, path string) (content []byte, err error) {
	content, err = ctx.GetFile(path)
	return
}

var (
	fileContent = []byte("Test123Test456Test789\n")
)
