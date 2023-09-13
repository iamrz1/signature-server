package config

import (
	"fmt"
	cerror "signature-server/error"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const EnvConfigFileKey = "CONFIG_FILE"
const DefaultConfigFileDir = "config.yaml"
const DefaultEnvironment = "dev"

// App ...
type App struct {
	Environment string
	ServerPort  int
	SystemPort  int
	Timeout     time.Duration
	DaemonKey   string
	LogLevel    string
}

func (cnf *App) validate() error {
	err := cerror.ValidationError{}
	if cnf.ServerPort == 0 {
		err.Add("server_port", "required")
	}
	if cnf.SystemPort == 0 {
		err.Add("system_port", "required")
	}
	if cnf.Timeout == 0 {
		err.Add("timeout", "required")
	}
	if cnf.DaemonKey == "" {
		err.Add("daemon_key", "required")
	}

	if len(err) > 0 {
		return fmt.Errorf("app configuration error: %v", err)
	}
	return nil
}

func (cnf *App) setDefaults() {
	if cnf.Environment == "" {
		cnf.Environment = DefaultEnvironment
	}
}

var appCnf App
var appErr error
var appOnce = sync.Once{}

func loadApp() {
	appCnf = App{
		Environment: viper.GetString("app.environment"),
		ServerPort:  viper.GetInt("app.server_port"),
		SystemPort:  viper.GetInt("app.system_port"),
		Timeout:     viper.GetDuration("app.timeout") * time.Second,
		DaemonKey:   viper.GetString("app.daemon_key"),
		LogLevel:    viper.GetString("app.log_level"),
	}
}

// AppCnf ...
func AppCnf() (App, error) {
	appOnce.Do(func() {
		read()
		loadApp()
		appCnf.setDefaults()
		appErr = appCnf.validate()
	})
	return appCnf, appErr
}
