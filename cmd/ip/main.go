package main

import (
	"flag"
	"log"

	"github.com/iofq/ip/http"
)

func main() {
	listen := flag.String("p", "8080", "Port to listen on")
	port := ":" + *listen

	server := http.New([]string{"X-Real-IP", "X-Forwarded-For"})
	log.Printf("Listening on http:%s", *listen)
	if err := server.ListenAndServe(port); err != nil {
		log.Fatal(err)
	} else {
	}
}
