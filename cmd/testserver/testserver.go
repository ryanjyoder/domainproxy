package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	message := "Welcome to the test server!"
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("remote:", r.RemoteAddr)
		fmt.Fprintln(w, message)
		//fmt.Println(sURL, "path:", r.URL.String())
	}))

	time.Sleep(1 * time.Hour)
}
