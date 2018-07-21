package domainproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type ProxyService interface {
	AddDomain(string, string) error
	DeleteDomain(string) error
	DomainExists(string) (bool, error)
}

type service struct {
	mux           *DomainMux
	domains       map[string]string
	domainsMutex  sync.Mutex
	domainStorage string
}

func NewProxyServer() (ProxyService, error) {
	s := &service{
		mux:           NewDomainMux(),
		domains:       make(map[string]string),
		domainsMutex:  sync.Mutex{},
		domainStorage: "domain_list.tab",
	}

	go http.ListenAndServe(":8080", s.mux)

	return s, nil
}

func (s *service) AddDomain(domain string, backend string) error {

	backendURL, err := url.Parse("http://" + backend)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	s.mux.SetHandler(domain, proxy)

	s.domainsMutex.Lock()
	s.domains[domain] = backend
	s.domainsMutex.Unlock()

	s.storeDomains()

	return nil
}

func (s *service) DeleteDomain(domain string) error {
	return nil
}

func (s *service) DomainExists(domain string) (bool, error) {
	s.domainsMutex.Lock()
	defer s.domainsMutex.Unlock()
	_, ok := s.domains[domain]
	return ok, nil
}

func (s *service) storeDomains() {

}
