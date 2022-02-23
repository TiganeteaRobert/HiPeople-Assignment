package main

import (
	"hipeople/image"

	"log"
	"net/http"
)

func main() {
	image.InitEndpoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
