package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	server.ListenAndServe()
}
