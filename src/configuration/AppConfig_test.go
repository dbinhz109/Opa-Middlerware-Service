package configuration

import (
	"testing"

	"github.com/spf13/viper"
)

/**

config.yml

configuration:
  profiles:
	  - dev

config.dev.yml

server:
  port: 8877

*/

func TestLoadingConfig(t *testing.T) {
	InitializeAppConfig()
	profiles := viper.GetStringSlice(CFG_CONFIGURATION_PROFILES)
	t.Logf("profiles: %v", profiles)
	if len(profiles) < 1 || profiles[0] != "dev" {
		t.Errorf("Expected profiles [dev]")
	}
	port := viper.GetInt(CFG_SERVER_PORT)
	if port != 8877 {
		t.Errorf("Expected dev profile port 8877 instead of %v", port)
	}

	// env APP_CONSUL_CONFIG_KEY
	consulKey := viper.GetString(CFG_CONSUL_CONFIGKEY)
	if consulKey != "appconfig" {
		t.Errorf("Expected ConsulKey appconfig instead of %v", consulKey)
	}
}
