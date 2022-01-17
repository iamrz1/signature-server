package config

import (
	"github.com/spf13/viper"
	"os"
	"signature-server/util"
	"sync"
)

var cnfOnce = sync.Once{}

func read() {
	cnfOnce.Do(func() {
		configFileDir := os.Getenv(EnvConfigFileKey)
		if configFileDir == "" {
			configFileDir = DefaultConfigFileDir
		}
		viper.SetConfigFile(configFileDir)
		err := viper.ReadInConfig()
		if err != nil {
			util.Fatalf("unable to read config file: %v", err)
		}
	})
}
