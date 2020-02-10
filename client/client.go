package hsfs

type Client interface {
	Auth(serverURL, pgpPrivKeyFilePath, basicAuthFilePath string) (Session, error)
	Post(path string, content []byte) error
	Get(path string) (content []byte, err error)
}
