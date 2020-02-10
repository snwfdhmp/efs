package services

import (
	"path/filepath"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	fs = afero.NewOsFs()
)

type Ctx interface {
	Log() *zap.Logger
	AES() AESService
	API() ApiService
	Auth() AuthenticationService
	FS() FSService

	PostFile(path string, content []byte) error
	GetFile(path string) ([]byte, error)
}

type ctx struct {
	logger                *zap.Logger
	aesService            AESService
	apiService            ApiService
	authenticationService AuthenticationService
	fsService             FSService
	fsPath                string
}

func NewCtx(logger *zap.Logger, fsPath string) (Ctx, error) {
	aesKey, err := afero.ReadFile(fs, filepath.Join(fsPath, "system/aes.key"))
	if err != nil {
		return nil, err
	}
	aesSvc := NewAESService(aesKey)
	fsSvc, err := NewFSService(aesSvc, fsPath)
	if err != nil {
		return nil, err
	}

	authSvc := NewAuthenticationService()

	ctx := ctx{
		logger:                logger,
		fsPath:                fsPath,
		fsService:             fsSvc,
		aesService:            aesSvc,
		authenticationService: authSvc,
	}

	api := NewAPIService(&ctx)
	ctx.apiService = api
	return &ctx, nil
}

func (c *ctx) Log() *zap.Logger {
	return c.logger
}

func (c *ctx) AES() AESService {
	return c.aesService
}

func (c *ctx) API() ApiService {
	return c.apiService
}

func (c *ctx) Auth() AuthenticationService {
	return c.authenticationService
}

func (c *ctx) FS() FSService {
	return c.fsService
}

func (c *ctx) PostFile(path string, content []byte) error {
	encrypted, err := c.AES().encrypt(content)
	if err != nil {
		return err
	}
	return c.FS().copyUserFile(path, encrypted)
}

func (c *ctx) GetFile(path string) ([]byte, error) {
	encrypted, err := c.FS().getUserFile(path)
	if err != nil {
		return nil, err
	}
	return c.AES().decrypt(encrypted)
}
