package main

import (
	"hipeople/store"

	"log"
	"net/http"
)

func main() {
	store.InitEndpoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
