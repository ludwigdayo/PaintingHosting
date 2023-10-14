SET BLOB_NAME=PaintingHosting
SET CGO_ENABLED=0

SET GOOS=windows
SET GOARCH=amd64
go build -o .\out\%BLOB_NAME%-%GOOS%-%GOARCH%.exe .

SET GOOS=darwin
SET GOARCH=amd64
go build -o .\out\%BLOB_NAME%-%GOOS%-%GOARCH% .

SET GOOS=linux
SET GOARCH=amd64
go build -o .\out\%BLOB_NAME%-%GOOS%-%GOARCH% .

