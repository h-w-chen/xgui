package web

import (
	"fmt"
	"net/http"
	"strconv"
)

var chY = make(chan int)

// GetYOnPort produces Y based on what is received from the tcp server
func GetYOnPort() int {
	return <-chY
}

// StartServer starts the listening HTTP server
func StartServer(port int) {
	http.HandleFunc("/y", func(w http.ResponseWriter, r *http.Request) {
		val := r.URL.Query().Get("val")
		y, err := strconv.Atoi(val)
		if err != nil {
			// discard the invalid y value
			return
		}
		chY <- y
	})

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
