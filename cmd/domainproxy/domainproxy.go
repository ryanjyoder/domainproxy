package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/ryanjyoder/domainproxy"
)

func main() {
	myService, err := domainproxy.NewProxyServer()
	if err != nil {
		log.Fatal("could not start service:", err)
	}
	server := http.Server{
		Handler: &myHandler{service: myService},
	}
	socketDir := os.Getenv("SNAP_DATA")
	fmt.Println("SNAP_DATA:", socketDir)
	if socketDir == "" {
		socketDir = "./"
	}
	socketPath := filepath.Join(socketDir, "go.sock")
	unixListener, err := net.Listen("unix", socketPath)
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
	service domainproxy.ProxyService
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
		domain := getDomainRegex.FindStringSubmatch(req.URL.String())[1]
		isSet, _ := h.service.DomainExists(domain)
		response, _ := json.Marshal(getResponse{IsSet: isSet})
		w.Write(response)
		return
	case req.Method == "PUT" && putDomainRegex.MatchString(req.URL.String()):
		args := putDomainRegex.FindStringSubmatch(req.URL.String())[1:]
		h.service.AddDomain(args[0], args[1])
		return
	default:

	}
}
