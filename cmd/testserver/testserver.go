package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	message := "Welcome to the test server!"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, message)
		//fmt.Println(sURL, "path:", r.URL.String())
	}))
	defer s.Close()

	fmt.Println("URL:", s.URL)
	time.Sleep(1 * time.Hour)
}
