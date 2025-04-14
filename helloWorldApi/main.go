package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	helloHandler := func(writer http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(writer, "hello world")
	}

	mux.HandleFunc("/", helloHandler)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
