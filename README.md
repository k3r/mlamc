# mlamc
# How to use it
```
# Get help
C:\>mlamc -h
Usage of mlamc:
  -d string
        Directories to analyze
  -disable-vt
        Disable Virus Total
  -e string
        Extensions to use for directory analysis (default ".exe;.dll;.js;.vbs")
  -f string
        Files to analyze
  -v    Verbose mode
 
 
```
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
