package main

import (
	"context"

	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/e2e/testdata"
)

func main() {
	config := config.NewConfig()
	testdata.StartServer(context.Background(), config)
}
