package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const defaultListen = ":8080"

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Recieved request: %s", b)
	fmt.Fprint(w, "ok")
}

func main() {
	http.HandleFunc("/", handler)

	var listen string
	if listen = os.Getenv("API_LISTEN"); listen == "" {
		listen = defaultListen
	}
	log.Printf("Starting http server on %s", listen)
	log.Fatal(http.ListenAndServe(defaultListen, nil))
}
