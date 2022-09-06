# gogs: Golang Game Server Framework

[![version](https://img.shields.io/github/v/tag/metagogs/gogs?label=version)](https://github.com/metagogs/gogs)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/metagogs/gogs/Go)
[![codecov](https://codecov.io/gh/metagogs/gogs/branch/main/graph/badge.svg?token=CLNQZ26H6X)](https://codecov.io/gh/metagogs/gogs)
[![goversion](https://img.shields.io/github/go-mod/go-version/metagogs/gogs)](https://github.com/metagogs/gogs)
[![license](https://img.shields.io/github/license/metagogs/gogs)](https://github.com/metagogs/gogs)
![GitHub issues](https://img.shields.io/github/issues/metagogs/gogs)
![GitHub last commit](https://img.shields.io/github/last-commit/metagogs/gogs)
![GitHub top language](https://img.shields.io/github/languages/top/metagogs/gogs)

gogs is an simple, fast and lightweight game server framewrok written in golang. It is designed to be easy to use and easy to extend. It will generate logic code from protobuf files, and you can use it to develop your game server. It is also a good choice for you to learn golang. It support websocket and webrtc datachannel.

---

### TODO
- [ ] Support metrics
- [ ] Support generate Unity C# SDK
- [ ] Support JS SDK
- [ ] Support Remote call
- [ ] Support tracing
- [ ] Add more examples
- [ ] Add more tests
- [ ] Add more documentation

## Getting Started
### Prerequisites
* [Go](https://golang.org/) >= 1.10
* [Protobuf](https://developers.google.com/protocol-buffers)
### Init your project
install the gogs
```
go install github.com/metagogs/gogs/tools/gogs@v0.0.9
```
init project
```
mkdir yourgame
cd yourgame
gogs init -p yourgame
```
edit your proto, add the game message, then generate the code
```
gogs go -f data.proto
```
run your game server
```
go mod tidy
go run main.go
```
### Generated Project
```
internal/
    logic/
        baseworld/
            bind_user_logic.go
    server/
        server.go
    svc/
        service_context.go
model/
    data.ep.go
    data.pb.go
config.yaml     
data.proto      
main.go
```
## Contributing
### Running the gogs tests
```
make test
```
This command will run both unit and e2e tests.
## License
[Apache License Version 2.0](./LICENSE)

