package image

import (
	"net/http"
)

func InitEndpoints() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer InfoLogger.Println(`-----------------`)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if r.Method == "GET" {
			getImage(w, r)
		} else if r.Method == "POST" {
			addImage(w, r)
		} else {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		}
	})
}
