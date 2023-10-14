SET BLOB_NAME=PaintingHosting
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build -o .\out\%BLOB_NAME%-%GOOS%-%GOARCH%.exe .

