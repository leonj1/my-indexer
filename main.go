package main

import (
	"log"
	"net/http"
	"my-indexer/router"
)

func main() {
	r := router.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
