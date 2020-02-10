package services

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FSService interface {
	getUserFile(path string) (encryptedFile []byte, err error) // get encrypted file corresponding to path
	copyUserFile(path string, content []byte) error            // create file in fs)
	// GetSystemFile(path string) (encryptedFile []byte, err error) // get encrypted file corresponding to path
	// CreateSystemFile(path string, content []byte) error          // create file in fs)
	getDir(name string) string
	getName(path string) string
}

type fsService struct {
	filenameSalt    []byte
	allowOverwrites bool
	basePath        string
}

func NewFSService(aes AESService, fsPath string) (FSService, error) {
	file, err := os.Open(filepath.Join(fsPath, "system/filename-salt.txt"))
	if err != nil {
		return nil, err
	}
	salt, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &fsService{
		filenameSalt:    salt,
		basePath:        fsPath,
		allowOverwrites: true,
	}, nil
}

func (l *fsService) getUserFile(path string) (encryptedFile []byte, err error) {
	filename := l.getName(path)
	file, err := os.OpenFile(filepath.Join(l.basePath, l.getDir(filename), filename), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

func (l *fsService) copyUserFile(path string, content []byte) error {
	filename := l.getName(path)
	filedir := l.getDir(filename)
	fpj := filepath.Join
	filepath := filepath.Join(l.basePath, filedir, filename)
	if _, err := os.Stat(filepath); err == nil {
		//verify
		if !l.allowOverwrites {
			log.Warn("overwrites not allowed ! cannot copy")
			return fmt.Errorf("fsService: allowOverwrites=false")
		}
	} else if os.IsNotExist(err) {
		if err = os.MkdirAll(fpj(l.basePath, filedir), 0700); err != nil {
			return err
		}
	} else {
		return err
	}

	log.Info(fmt.Sprintf("Writing %d bytes to %s", len(content), filepath))
	return ioutil.WriteFile(filepath, content, 0600)
}

func (l *fsService) getName(path string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(path+string(l.filenameSalt))))
}
func (l *fsService) getDir(name string) string {
	return fmt.Sprintf("files/%s/%s", name[0:2], name[2:4])
}
