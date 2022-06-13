package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/brensch/recreation"
)

func main() {

	// proxy.Director = recreation.Director

	// http.HandleFunc("/", HelloServer)
	http.HandleFunc("/", recreation.HandleAvailabilitiesSync)

	http.ListenAndServe(":8081", nil)

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----")

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(dump))

}
