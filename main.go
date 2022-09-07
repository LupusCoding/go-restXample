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
	mux.HandleFunc("/user/list", userList)

	fmt.Println("Server started. Listening on Port: " + strconv.FormatUint(uint64(port), 10))
	err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), mux)

	fmt.Print(err)
}

func logRequest(r *http.Request) {
	fmt.Println("DEBUG: Received " + r.Method + " request:")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	transfer.ParseUnsigned(r, struct{}{})
	// Do health checks
	transfer.RespondUnsigned(w, struct{}{}, true)
}

func userList(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// create data structure
	type reqData struct {
		Amount int
	}
	type resData struct {
		Amount int
		Users  []string
	}

	// handle request
	var rqd reqData
	req, err := transfer.ParseSigned(r, rqd, "./keys/public_key.pem")
	if err != nil {
		panic(err)
	}

	// handle response
	var rsd resData
	reqDatePointer := req.Data.(reqData)
	rsd.Amount = reqDatePointer.Amount
	for len(rsd.Users) < rsd.Amount {
		rsd.Users = append(rsd.Users, "Max Musermann")
	}
	err = transfer.RespondSigned(w, rsd, "./keys/private_key.pem")
	if err != nil {
		panic(err)
	}
}
