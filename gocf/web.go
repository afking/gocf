package gocf

import (
	"fmt"
	"log"
	"net/http"
)

func web() {
	port := "8080"

	http.HandleFunc("/", webServer)

	log.Println("Running on Port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("error handling %q: %v", r.RequestURI, err)
		}
	}
}

func webServer(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello crazyflie")
	return nil
}
