package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	efsRootDir = os.Getenv("EFS_ROOT_DIT")
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "create-fs",
				Usage: "create filesystem root dir",
				Action: func(c *cli.Context) error {
					fsName := c.Args().Get(0)
					if fsName == "" {
						fmt.Println("give fsName")
					}
					if err := os.MkdirAll(filepath.Join(efsRootDir, fsName, "system"), 0700); err != nil {
						return err
					}

					if err := os.MkdirAll(filepath.Join(efsRootDir, fsName, "system", "users"), 0700); err != nil {
						return err
					}

					if err := os.MkdirAll(filepath.Join(efsRootDir, fsName, "files"), 0700); err != nil {
						return err
					}

					// aes key
					aesBuf := make([]byte, 32)
					_, err := rand.Read(aesBuf)
					if err != nil {
						return err
					}

					aesB64Buf := bytes.NewBuffer([]byte{})
					_, err = base64.NewEncoder(base64.StdEncoding, aesB64Buf).Write(aesBuf)
					if err != nil {
						return err
					}

					err = ioutil.WriteFile(filepath.Join(efsRootDir, fsName, "system", "aes.key"), aesB64Buf.Bytes(), 0600)
					if err != nil {
						return err
					}

					saltBuf := make([]byte, 64)
					_, err = rand.Read(saltBuf)
					if err != nil {
						return err
					}

					err = ioutil.WriteFile(filepath.Join(efsRootDir, fsName, "system", "filename-salt.txt"), []byte(fmt.Sprintf("%x", saltBuf)), 0600)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
