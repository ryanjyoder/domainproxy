package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/ryanjyoder/domainproxy"
)

func main() {
	mux := domainproxy.NewDomainMux()
	loadStoredDomains(mux)
	go http.ListenAndServe(":8080", mux)

	server := http.Server{
		Handler: &myHandler{mux: mux},
	}

	unixListener, err := net.Listen("unix", "go.sock")
	// Create the socket to listen on:

	if err != nil {
		log.Fatal(err)
		return
	}

	// Unix sockets must be unlink()ed before being reused again.

	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		// Stop listening (and unlink the socket if unix type):
		unixListener.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)
	if err != nil {
		panic(err)
	}
	server.Serve(unixListener)
}

type myHandler struct {
	mux *domainproxy.DomainMux
}

var getDomainRegex = regexp.MustCompile("/domain/([^/]+)$")
var putDomainRegex = regexp.MustCompile("/domain/([^/]+)/([^/]+)$")

type getResponse struct {
	IsSet bool `json:"is_set"`
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("request:", req.URL.String())
	switch {
	case req.Method == "GET" && getDomainRegex.MatchString(req.URL.String()):
		isSet := h.mux.IsDomainSet(getDomainRegex.FindStringSubmatch(req.URL.String())[1])
		response, _ := json.Marshal(getResponse{IsSet: isSet})
		w.Write(response)
		return
	case req.Method == "PUT" && putDomainRegex.MatchString(req.URL.String()):
		args := putDomainRegex.FindStringSubmatch(req.URL.String())[1:]
		backendURL, _ := url.Parse(args[0])
		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		domain := args[0]
		fmt.Println("setting:", domain, backendURL.String())
		h.mux.SetHandler(domain, proxy)
	default:

	}
}

func loadStoredDomains(mux *domainproxy.DomainMux) {
	// first the "instance" should be launched
	tekURL, _, _ := domainproxy.LaunchTestServer("Welcome to TekCitadel Server!")

	//set up proxy
	tekProxy := httputil.NewSingleHostReverseProxy(tekURL)

	tekDomain := "tekcitadel.com"
	mux.SetHandler(tekDomain, tekProxy)

	url, _ := url.Parse("http://10.111.74.111:8080")
	mux.SetHandler("shell.ryanjyoder.com", httputil.NewSingleHostReverseProxy(url))
}
