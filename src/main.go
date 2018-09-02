//API server adapted from https://pragmacoders.com/building-a-json-api-in-golang/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	Populate()

	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/query", handleQryMessage).Methods("GET")
	//router.HandleFunc("/m/{msg}", handleURLMessage).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func handleQryMessage(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	//message := vars.Get("msg")
	person := vars.Get("Person")
	relation := vars.Get("Relation")

	fmt.Printf("Searching for %s %s\n", person, relation)
	relatives := search(person, relation)
	if relatives == nil || len(relatives) == 0 {
		relatives = append(relatives, "none")
	}

	json.NewEncoder(w).Encode(map[string][]string{"message": relatives})

	//json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleURLMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	message := vars["msg"]

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
