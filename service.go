package domainproxy

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

func NewProxyServer(domainListFile string) (ProxyService, error) {
	s := &service{
		mux:           NewDomainMux(),
		domains:       make(map[string]string),
		domainsMutex:  sync.Mutex{},
		domainStorage: domainListFile,
	}

	s.loadFromStorage()

	go http.ListenAndServe(":8080", s.mux)

	return s, nil
}

func (s *service) AddDomain(domain string, backend string) error {

	backendURL, err := url.Parse(backend)
	if err != nil {
		return err
	}

	s.domainsMutex.Lock()
	currentBackend, ok := s.domains[domain]
	s.domainsMutex.Unlock()
	if ok && currentBackend == backend { // nothing to do
		return nil
	}
	if ok { // changing the backend delete the current
		s.DeleteDomain(domain)
	}

	s.domainsMutex.Lock()
	s.domains[domain] = backend
	s.domainsMutex.Unlock()

	s.storeDomains()

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	s.mux.SetHandler(domain, proxy)

	return nil
}

func (s *service) DeleteDomain(domain string) error {
	s.domainsMutex.Lock()
	delete(s.domains, domain)
	s.domainsMutex.Unlock()
	return nil
}

func (s *service) DomainExists(domain string) (bool, error) {
	s.domainsMutex.Lock()
	defer s.domainsMutex.Unlock()
	_, ok := s.domains[domain]
	return ok, nil
}

func (s *service) loadFromStorage() error {
	fileio, err := os.Open(s.domainStorage)
	defer fileio.Close()
	if err != nil {
		return err
	}
	reader := csv.NewReader(fileio)
	reader.Comma = '\t'
	domains, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i := range domains {
		if len(domains[i]) < 2 {
			return fmt.Errorf("expecting tab: %v", domains[i])
		}
		s.AddDomain(domains[i][0], domains[i][1])
	}

	return nil
}

func (s *service) storeDomains() error {
	records := make([][]string, len(s.domains))
	s.domainsMutex.Lock()
	i := 0
	for domain, backend := range s.domains {
		records[i] = []string{domain, backend}
	}
	s.domainsMutex.Unlock()

	fileio, err := os.Create(s.domainStorage)
	defer fileio.Close()
	if err != nil {
		return err
	}
	writer := csv.NewWriter(fileio)
	writer.Comma = '\t'
	writer.WriteAll(records)

	return nil
}
