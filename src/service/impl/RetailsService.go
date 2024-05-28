package impl

import (
	"github.com/spf13/viper"
)

type RetailsService struct {
	AIABaseUrl string
}

func NewRetailsService() *RetailsService {
	svc := &RetailsService{AIABaseUrl: viper.GetString("services.aia.baseUrl")}
	return svc
}
