package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Temporary struct to create Pokemon objects before adding SQL database
type Pokemon struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

// Added log.Fatal so that it sends an error if the server crashes
func main() {
	http.HandleFunc("/pokemon", handlePokemon)
	http.HandleFunc("/pokemon/", handlePokemon) // Emulates "/pokemon/{id}" on a framework as this is a limitation of only using net/http
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sorry this method is not supported"))
	}
}

func handleGetPokemon(w http.ResponseWriter, r *http.Request) error {
	// Checks to see if the GET request was just on "/pokemon" If not, it removes the id from the end to get a specific pokemon
	if r.URL.Path == "/pokemon" {
		pokemon := Pokemon{
			Number: 25,
			Name:   "Pikachu",
		}
		return WriteJSON(w, http.StatusOK, pokemon)
	}

	// Testing to show that I can trim the id from the URL. Using to WriteJSON to write the id response back to the client.
	id := strings.TrimPrefix(r.URL.Path, "/pokemon/")
	return WriteJSON(w, http.StatusOK, id)
}
