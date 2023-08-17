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
	http.HandleFunc("/pokemon", handlePokemon)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Why do I need to put writeHeader() underneath Header.Set()??? If I don't then it sends as plain text, not JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Hard coded in a pokemon so I could get JSON encoding working before adding SQL
func handlePokemon(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPokemon(w, r)
	case "POST":
		w.Write([]byte("This is a post request"))
	case "PATCH":
		w.Write([]byte("This is a patch request"))
	case "DELETE":
		w.Write([]byte("This is a delete request"))
	default:
		w.Write([]byte("Sorry this method is not supported"))
	}
}

func handleGetPokemon(w http.ResponseWriter, r *http.Request) error {
	pokemon := Pokemon{
		Number: 25,
		Name:   "Pikachu",
	}
	return WriteJSON(w, http.StatusOK, pokemon)
}
