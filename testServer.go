package domainproxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func LaunchTestServer(message string) (sURL *url.URL, stopChan chan int, err error) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, message)
		//fmt.Println(sURL, "path:", r.URL.String())
	}))
	sURL, err = url.Parse(s.URL)

	stopChan = make(chan int)
	go func() {
		<-stopChan
		s.Close()
	}()
	return
}
