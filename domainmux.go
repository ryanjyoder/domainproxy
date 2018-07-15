package domainproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

type DomainMux struct {
	handlers       map[string]http.Handler
	defaultHandler http.Handler
	handlerMutex   sync.RWMutex
}

func (h *DomainMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := (strings.Split(r.Host, ":"))[0]
	fmt.Println("host:", domain)
	if handler, ok := h.handlers[domain]; ok {
		handler.ServeHTTP(w, r)
	} else {
		h.defaultHandler.ServeHTTP(w, r)
	}
}

func NewDomainMux() *DomainMux {
	mux := &DomainMux{}
	mux.handlers = make(map[string]http.Handler)
	defaultURL, _, _ := LaunchTestServer("default Server here!")
	defaultProxy := httputil.NewSingleHostReverseProxy(defaultURL)
	mux.defaultHandler = defaultProxy
	return mux
}

func (mux *DomainMux) SetHandler(domain string, handler http.Handler) {
	mux.handlerMutex.Lock()
	mux.handlers[domain] = handler
	mux.handlerMutex.Unlock()
}
func (mux *DomainMux) IsDomainSet(domain string) bool {
	mux.handlerMutex.RLock()
	_, ok := mux.handlers[domain]
	mux.handlerMutex.RUnlock()
	return ok
}
func (mux *DomainMux) DeleteHandler(domain string) {
	mux.handlerMutex.RLock()
	_, ok := mux.handlers[domain]
	mux.handlerMutex.RUnlock()
	if ok {
		mux.handlerMutex.Lock()
		delete(mux.handlers, domain)
		mux.handlerMutex.Unlock()
	}
}
