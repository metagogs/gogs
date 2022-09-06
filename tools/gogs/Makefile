version := $(shell /bin/date "+%Y-%m-%d %H:%M")

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o gogs main.go
	$(if $(shell command -v upx), upx gogs)
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o gogs main.go
	$(if $(shell command -v upx), upx gogs)
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o gogs.exe main.go
	$(if $(shell command -v upx), upx gogs.exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o gogs main.go
	$(if $(shell command -v upx), upx gogs)
