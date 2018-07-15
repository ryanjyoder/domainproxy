package domainproxy

import (
	"sync"
)

type ProxyService interface {
	//AddDomain(string, string) error
	//DeleteDomain(string, string) error
}

type config struct {
	configFile string
}

type service struct {
	mux          *DomainMux
	domains      map[string]string
	domainsMutex sync.Mutex
}

func NewProxyServer(configs config) (ProxyService, error) {
	s := &service{
		mux:          NewDomainMux(),
		domains:      make(map[string]string),
		domainsMutex: sync.Mutex{},
	}

	return s, nil
}
