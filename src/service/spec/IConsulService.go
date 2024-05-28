package spec

import (
	"github.com/hashicorp/consul/api"
)

type IConsulService interface {
	GetService(serviceName string) ([]*api.ServiceEntry, error)
	RegisterService(serviceName, serviceId, serviceAddress, servicePort string) error
	GetConfig(key string) (string, error)
}
