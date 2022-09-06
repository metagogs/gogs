#!/bin/bash

go build -o gogs main.go
chmod +x gogs
cp gogs `go env GOPATH`/bin
rm gogs

