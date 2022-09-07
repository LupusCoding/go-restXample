package main

import (
	"fmt"
	"net/http"
	"restXample/transfer"
	"strconv"
)

var port uint = 8080

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health/check", healthCheckHandler)

	fmt.Println("Server started. Listening on Port: " + strconv.FormatUint(uint64(port), 10))
	err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), mux)

	fmt.Print(err)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DEBUG: Received " + r.Method + " request:")

	transfer.ParseUnsigned(r, struct{}{})
	// Do health checks
	transfer.RespondUnsigned(w, struct{}{}, true)
}
