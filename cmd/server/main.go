package main

import (
	"github.com/antik9/aws-web-app/internal/server"
)

func main() {
	http.HandleFunc("/", server.HelloWorld())
	log.Fatal(
		http.ListenAndServe(":8000", nil),
	)
}
