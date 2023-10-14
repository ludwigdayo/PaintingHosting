BLOB_NAME=PaintingHosting
CGO_ENABLED=1
GOOS=linux
GOARCH=amd64
go build -o ./out/${BLOB_NAME}-${GOOS}-${GOARCH} .