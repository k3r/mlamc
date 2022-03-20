SET FILES=mlamc.go api.go

SET GOOS=darwin
SET GOARCH=amd64
go build -o mlamc.mac %FILES%

SET GOOS=windows
SET GOARCH=amd64
go build -o mlamc.exe %FILES%

SET GOOS=linux
SET GOARCH=amd64
go build -o mlamc.linux %FILES%