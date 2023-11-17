# gogs: Golang Game Server Framework


[![version](https://img.shields.io/github/v/tag/metagogs/gogs?label=version)](https://github.com/metagogs/gogs)
[![codecov](https://codecov.io/gh/metagogs/gogs/branch/main/graph/badge.svg?token=CLNQZ26H6X)](https://codecov.io/gh/metagogs/gogs)
![GitHub issues](https://img.shields.io/github/issues/metagogs/gogs)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/metagogs/gogs/Go)
![GitHub Release Date](https://img.shields.io/github/release-date/metagogs/gogs)
![GitHub last commit](https://img.shields.io/github/last-commit/metagogs/gogs)
[![goversion](https://img.shields.io/github/go-mod/go-version/metagogs/gogs)](https://github.com/metagogs/gogs)
[![license](https://img.shields.io/github/license/metagogs/gogs)](https://github.com/metagogs/gogs)
![GitHub top language](https://img.shields.io/github/languages/top/metagogs/gogs)



gogs is a simple, fast and lightweight game server framework written in golang. It is designed to be easy to use and easy to extend. It will generate logic code from protobuf files, and you can use it to develop your game server. It's also a good starting point for you to learn golang. It supports websocket and webrtc datachannel.

[Untiy Meta City Demo Online](https://metagogs.github.io/metacity/)

---

## TODO
- [ ] Support metrics
- [x] Support generate Unity C# SDK
- [ ] Support generate JS SDK
- [ ] Support generate Golang SDK
- [ ] Support Remote call
- [ ] Support tracing
- [x] Support gogs generate docker file
- [x] Support gogs generate k8s yaml
- [ ] Support custom game packet protocol
- [ ] Support kubegame controller, create game pod with api
- [ ] Add more examples
- [ ] Add more tests
- [ ] Add more documentation
- [ ] Test coverage reaches 80% 
- [ ] k8s friendly, hot reload?


## Getting Started
### Prerequisites
* [Go](https://golang.org/) >= 1.21
* [Protobuf](https://developers.google.com/protocol-buffers)
### Init your project
install the gogs
```
go install github.com/metagogs/gogs/tools/gogs@v0.2.4
```
init project
```
mkdir yourgame
cd yourgame
gogs init -p yourgame

Flags:
 -p your go package name
```
edit your proto, add the game message, then generate the code
```
gogs go -f data.proto

Flags:
 -f proto file path
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

### Generated Unity C# Code
this will generate a unity code, you can use it to test your game server. And you need use the [Unity Protobuf](./unity) to run the code.
```
 gogs csharp -f data.proto

 Flags:
 -f proto file path
 -g generate unity c# gogs library, you should use it when you generate code first time
```

### Deployment
```
gogo docker // generate Dockerfile

go run main --deployment // generate k8s yaml
--svc use service to expose your game server not hostport
--name your game server name
--namepsace your k8s namespace
```

## How to encode/decode the message
### Packet Protocol
gogs uses 8 bytes as the protocol header

protocol header = flag + version + action index + message data length
```
0     byte  flag         8bit    protocol flag, awalys 0x7E 
1     byte  version      5bit    protocel version
            encodeType   3bit    protocel encode type
2     byte  packetType   2bit    packet message type 1 system 2 server
            component    6bit    message component index
3.4   byte  action       16bit   message action index
5.6.7 byte  length       24bit   message length
```

### What is the action index

action index = packetType + component index + action index
```protobuf
// @gogs:Components
message Components {
    BaseWorld BaseWorld = 1; // 1 is the component index
}

message BaseWorld {
    BindUser BindUser = 1; // 1 is the action index

    BindSuccess BindSuccess = 2; // 2 is the action index
}
message BindUser {
    string uid = 1;
}


// @gogs:ServerMessage
message BindSuccess {
}
```

like this proto, the BindUser and BindSuccess is the message comunication between client and server

BindUser action index = packetType <<22 | component <<16 | action = 2 << 22 | 1 << 16 | 1 = 0x810001

BindSuccess action index = packetType <<22 | component <<16 | action = 2 << 22 | 1 << 16 | 2 = 0x810002

### Packet encode & decode
gogs has three encode&decode types - encodeType in protocol header
- 0 json encode&decode without protocol header
- 1 json encode&decode with protocol header
- 2 protobuf encode&decode with protocol header

### Packet message

**message with encode type 0 (json without protocol header)**

message = **JSON binary data**

`gogs retrieves the action index from the message, then gets the filed type and decodes the message, finally it calls the logic function. The json message without protocol header should add a filed named action, the value is the filed name`

```json
{
	"action": "BindUser",
	"uid": "123"
}
```

```golang
app.UseDefaultEncodeJSONWithHeader()
```

---
**message with encode type 1 (json with protocol header)**

message = **8 bytes protocol header** + **JSON binary data**

```golang
app.UseDefaultEncodeJSON()
```
---
**message with encode type 2 (protobuf with protocol header)**

message = **8 bytes protocol header** + **protobuf binary data**

```golang
app.UseDefaultEncodeProto()
```




## Contributing
### Running the gogs tests
```
make test
```
This command will run both unit and e2e tests.


## Demo
  + [Base demo generated by gogs](./examples/basedemo)
  + [Untiy Meta City Demo](https://github.com/metagogs/metacity)
  + [Untiy Meta City Demo Online](https://metagogs.github.io/metacity/)
  + [gtc - terminal chat app](https://github.com/szpnygo/gtc)
  
## License
[Apache License Version 2.0](./LICENSE)

