package gocf

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

func web() {
	port := "4000"

	http.Handle("/", websocket.Handler(errorHandler(socketHandler)))

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

func server(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello crazyflie")
	return nil
}

func socketHandler(c *websocket) {
	var s string
	fmt.Fscan(c, &s)
	fmt.Println("Received: ", s)
	fmt.Fprint(c, "How do you do?")
}
