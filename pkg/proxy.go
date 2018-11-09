package pkg

import (
	"fmt"
	"log"
)

type ProxyInfo struct {
	DefaultKey string
	Protocol   string
	User       string
	Password   string
	Endpoints  map[string]string
}

func (i ProxyInfo) GenerateProxy(key string) string {
	if _, ok := i.Endpoints[key]; !ok {
		log.Printf("proxy key <%s> not found, using default <%s>\n", key, i.DefaultKey)
		key = i.DefaultKey
	}

	if endpoint, ok := i.Endpoints[key]; ok {
		return fmt.Sprintf("--proxy-server=\"%s://%s:%s@%s\"", i.Protocol, i.User, i.Password, endpoint)
	}

	return ""
}
