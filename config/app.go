package config

import (
	"fmt"
	"github.com/spf13/viper"
	cerror "signature-server/error"
	"sync"
	"time"
)

const EnvConfigFileKey = "CONFIG_FILE"
const DefaultConfigFileDir = "config.yml"

// App ...
type App struct {
	ServerPort int
	SystemPort int
	Timeout    time.Duration
	DaemonKey  string
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

var appCnf App
var appErr error
var appOnce = sync.Once{}

func loadApp() {
	appCnf = App{
		ServerPort: viper.GetInt("app.server_port"),
		SystemPort: viper.GetInt("app.system_port"),
		Timeout:    viper.GetDuration("app.timeout") * time.Second,
		DaemonKey:  viper.GetString("app.daemon_key"),
	}
}

// AppCnf ...
func AppCnf() (App, error) {
	appOnce.Do(func() {
		read()
		loadApp()
		appErr = appCnf.validate()
	})
	return appCnf, appErr
}
