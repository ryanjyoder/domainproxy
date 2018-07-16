package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"bitbucket.org/ryanjyoder/domainproxy"
)

func main() {
	mux := domainproxy.NewDomainMux()
	go addDomains(mux)
	http.ListenAndServe(":8080", mux)
}

func addDomains(mux *domainproxy.DomainMux) {
	// first the "instance" should be launched
	tekURL, _, _ := domainproxy.LaunchTestServer("Welcome to TekCitadel Server!")

	//set up proxy
	tekProxy := httputil.NewSingleHostReverseProxy(tekURL)

	tekDomain := "tekcitadel.com"
	mux.SetHandler(tekDomain, tekProxy)

	url, _ := url.Parse("http://10.111.74.111:8080")
	mux.SetHandler("shell.ryanjyoder.com", httputil.NewSingleHostReverseProxy(url))
}
