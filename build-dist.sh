mkdir -p dist
cd ./server
go get ./cmd/efsctl/...
go build -o ../dist/efsctl ./cmd/efsctl