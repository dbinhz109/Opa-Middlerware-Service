package impl

import (
	"errors"
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type consulService struct {
	client *consulapi.Client
}

func NewConsulService() *consulService {
	consulAddr := viper.GetString("services.consul.server")
	if consulAddr == "" {
		return nil
	}
	client, err := consulapi.NewClient(&consulapi.Config{Address: consulAddr})
	if err != nil {
		log.Println(err)
		return nil
	}
	return &consulService{
		client: client,
	}
}

// GetService gets a service from consul by name
func (s *consulService) GetService(serviceName string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := s.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// RegisterService registers a service to consul
func (s *consulService) RegisterService(serviceName, serviceId, serviceAddress, servicePort string) error {
	if s == nil || s.client == nil {
		return errors.New("consul client is nil")
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = serviceId
	registration.Name = serviceName
	registration.Address = serviceAddress
	registration.Port, _ = strconv.Atoi(servicePort)
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           "http://" + serviceAddress + ":" + servicePort + "/healthz",
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	err := s.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

// Get config from consul
func (s *consulService) GetConfig(key string) (string, error) {
	kv, _, err := s.client.KV().Get(key, nil)
	if err != nil {
		return "", err
	}
	return string(kv.Value), nil
}
