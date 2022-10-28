package Utils

import (
	"fmt"
	"strings"
)

type ServiceCredentialType struct {
	Username string
	Password string
}

type Config struct {
	//"PROD"/"DEV"
	EnvType           string
	URLMapper         map[string]string
	ServiceCredential map[string]ServiceCredentialType
}

func InitConfig() Config {
	servicesNames := []string{"rabbitmq", "customerservice", "productservice"}
	servicesCredentials := []string{"rabbitmq"}

	newConfig := Config{
		EnvType:           GetEnv("ENVTYPE", "DEV"),
		URLMapper:         make(map[string]string),
		ServiceCredential: make(map[string]ServiceCredentialType),
	}

	for _, v := range servicesNames {
		newConfig.URLMapper[v] = GetEnv(fmt.Sprintf("%s%s", strings.ToUpper(v), "_URL"), "localhost")
	}

	for _, v := range servicesCredentials {
		newConfig.ServiceCredential[v] = ServiceCredentialType{
			Username: GetEnv(fmt.Sprintf("%s%s", strings.ToUpper(v), "_USERNAME"), "username"),
			Password: GetEnv(fmt.Sprintf("%s%s", strings.ToUpper(v), "_PASSWORD"), "password"),
		}
	}

	return newConfig
}
