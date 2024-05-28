package configuration

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-app/src/logger"
	"go-app/src/service"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

const (
	CFG_APPLICATION_NAME                   = "application.name"
	CFG_SERVER_PORT                        = "server.port"
	CFG_SERVER_HOST                        = "server.host"
	CFG_CONFIGURATION_PROFILES             = "configuration.profiles"
	CFG_ELASTICSEARCH_URL                  = "elasticsearch.url"
	CFG_APM_BROKERS                        = "apm.brokers"
	CFG_BASIC_USER                         = "auth.basic.user"
	CFG_BASIC_PASSWORD                     = "auth.basic.password"
	CFG_DATABASE_MONGODB_CONNECTION_STRING = "database.mongodb.connection-string"
	CFG_KAFKA_BROKERS                      = "kafka.brokers"
	CFG_CONSUL_SERVER                      = "services.consul.server"
	CFG_CONSUL_CONFIGKEY                   = "services.consul.config_key"
)

// InitializeAppConfig read and set default config values
func InitializeAppConfig() {
	logger.Info("Loading configuration")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.Set("server.ip", getOutboundIP())

	viper.SetDefault(CFG_APPLICATION_NAME, "GoApp")
	viper.SetDefault(CFG_SERVER_PORT, "8080")
	viper.SetDefault(CFG_SERVER_HOST, "")
	viper.SetDefault(CFG_CONFIGURATION_PROFILES, []string{})

	viper.SetDefault(CFG_CONSUL_SERVER, "http://localhost:8500")
	viper.SetDefault(CFG_CONSUL_CONFIGKEY, "goapp.configuration.json")

	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}
	profiles := viper.GetStringSlice(CFG_CONFIGURATION_PROFILES)
	for _, p := range profiles {
		profileCfg := "config." + p
		viper.SetConfigName(profileCfg)
		fmt.Printf("Loading %v\n", profileCfg)
		err := viper.MergeInConfig()
		if err != nil {
			logger.Info("MergeInConfig", zap.Error(err))
		}
	}

	consulServer := viper.GetString(CFG_CONSUL_SERVER)
	consulConfigKey := viper.GetString(CFG_CONSUL_CONFIGKEY)
	if consulServer != "" && consulConfigKey != "" {
		logger.Info("Loading configuration from Consul")
		viper.AddRemoteProvider("consul", consulServer, consulConfigKey)
		viper.SetConfigType("yaml")
		if err := viper.ReadRemoteConfig(); err != nil {
			logger.Info("ReadRemoteConfig", zap.Error(err))
		}
	}

	// GetConsulService("lsmanagement")

	randomBytes := make([]byte, 6)
	nBytes, err := io.ReadAtLeast(rand.Reader, randomBytes, 4)
	if err != nil || nBytes != 6 {
		logger.Info("NotEnoughRandom", zap.Error(err))
	}
	if os.Getenv("APP_APPLICATION_INSTANCEID") == "" {
		appInstanceId := hex.EncodeToString(randomBytes)
		// viper.Set("application.instanceId", appInstanceId)
		err = os.Setenv("APP_APPLICATION_INSTANCEID", appInstanceId)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Setenv", zap.Any("Setenv", appInstanceId))
	}
	logger.InitLoggingData()
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	logger.Info("ip", zap.Any("ip", localAddr.IP))
	return localAddr.IP
}

type jmap = service.JsonMap

func GetConsulService(serviceName string) {
	defer recover()
	consulServer := viper.GetString(CFG_CONSUL_SERVER)
	url := "http://" + consulServer + "/v1/catalog/service/" + serviceName
	resp, err := http.Get(url)
	if err != nil {
		logger.Info("Connection refused", zap.Error(err))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Empty body", zap.Error(err))
		return
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Info("result", zap.Error(err))
		return
	}
	logger.Info("resultbody", zap.Any("resultbody", result))
	var n int64
	if len(result) == 0 {
		logger.Info("Empty body", zap.Error(err))
		return
	}
	if int64(len(result)-1) <= 0 {
		n = 0
	} else {
		data, err := rand.Int(rand.Reader, big.NewInt(int64(len(result)-1)))
		if err != nil {
			panic(err)
		}
		n = data.Int64()
	}
	ipAddressConsul := result[n]["ServiceTaggedAddresses"].(jmap)["wan_ipv4"].(jmap)["Address"]
	ipAddressResp := ipAddressConsul.(string)
	portConsul := result[n]["ServiceTaggedAddresses"].(jmap)["wan_ipv4"].(jmap)["Port"]
	portResp := portConsul.(float64)
	dest := "http://" + ipAddressResp + ":" + strconv.FormatFloat(portResp, 'f', -1, 64)
	viper.Set("services.management.baseUrl", dest)
}
