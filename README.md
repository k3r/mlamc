# mlamc

## How to build
On Windows in the mlamc git clone:
```
SET FILES=mlamc.go api.go

SET GOOS=darwin
SET GOARCH=amd64
go build %FILES%

SET GOOS=windows
SET GOARCH=amd64
go build %FILES%

SET GOOS=linux
SET GOARCH=amd64
go build %FILES%
```
