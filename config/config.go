package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug       bool // debug mode
	Gopprof     bool // start pprof
	GopprofAddr int  // the pprof server http port
	AdminPort   int  // admin server http port and used for health check

	AgentMessageBufferSize int  // agent message buffer size
	AgentHeartBeatTimeout  int  // agent heartbeat timeout, second
	AgentHeartBeatLog      bool // show the agent heartbeat log

	SendMessageLog    bool // show the send message log
	ReceiveMessageLog bool // show the receive message log

	StaredCallback func() // just for testing
}

func NewDefaultConfig() *Config {
	return &Config{
		Debug:       false,
		Gopprof:     false,
		GopprofAddr: 9998,
		AdminPort:   9999,

		AgentMessageBufferSize: 1000,
		AgentHeartBeatTimeout:  20,
		AgentHeartBeatLog:      false,

		SendMessageLog:    false,
		ReceiveMessageLog: false,
	}
}

func setDefaultConfig() {
	config := NewDefaultConfig()
	defaultsMap := map[string]interface{}{
		"debug":                  config.Debug,
		"gopprof":                config.Gopprof,
		"gopprofaddr":            config.GopprofAddr,
		"adminport":              config.AdminPort,
		"agentmessagebuffersize": config.AgentMessageBufferSize,
		"agentheartbeattimeout":  config.AgentHeartBeatTimeout,
		"agentheartbeatlog":      config.AgentHeartBeatLog,
		"sendmessagelog":         config.SendMessageLog,
		"receivemessagelog":      config.ReceiveMessageLog,
	}

	for param := range defaultsMap {
		if viper.Get(param) == nil {
			viper.SetDefault(param, defaultsMap[param])
		}
	}
}

func NewConfig(file ...string) *Config {
	if len(file) > 0 {
		viper.SetConfigFile(file[0])
	} else {
		viper.SetConfigFile("config.yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetEnvPrefix("GOGS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	setDefaultConfig()
	viper.AutomaticEnv()

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return config
}
