package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	started := time.Now()

	if err := connectDatabase(); err != nil {
		log.Fatal(err)
	}
	defer databaseConn.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		duration := time.Now().Sub(started)

		if duration.Seconds() > 1000 {
			log.Println("Timeout triggered")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`Im timed out`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`Hello gopher`))
		}
	})

	http.HandleFunc("/aligator", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Mr Aligator")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
