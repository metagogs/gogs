module {{.ProjectPackage}}

{{.GoVersion}}

require (
	github.com/metagogs/gogs {{.GoGSVersion}}
	go.uber.org/zap v1.23.0
	google.golang.org/protobuf v1.28.1
)