package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Temporary struct to create Pokemon objects before adding SQL database
type Pokemon struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

// Added log.Fatal so that it sends an error if the server crashes
func main() {
	http.HandleFunc("/GetPokemon", handleGetPokemon)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Hard coded in a pokemon so I could get JSON encoding working before adding SQL
func handleGetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pokemon := Pokemon{
		Number: 025,
		Name:   "Pikachu",
	}

	json.NewEncoder(w).Encode(pokemon)
}
