package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, World!")
		}),
	}
	fmt.Println("Starting server on :8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
