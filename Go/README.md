



```text
win build linux

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

win build mac

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build

mac build linux and windows x86_64

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

linux build mac and windows x86_64

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```
