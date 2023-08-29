package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

// struct to create Pokemon objects before encoding into JSON
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

// Encodes data into JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Handles the incoming http requests and routes them depending on the method used.
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

func handleGetPokemon(w http.ResponseWriter, r *http.Request) any {
	var allPokemon []Pokemon
	var query string

	// Opens the pokemon database & defers closing until the end of the function
	db, errs := sql.Open("sqlite", "./pokemon.db")
	if errs != nil {
		log.Fatal(errs)
	}
	defer db.Close()

	// If the client added an id to the end of the URL this section extracts the id
	id := strings.TrimPrefix(r.URL.Path, "/pokemon/")

	// Checks to see if the GET request was just on "/pokemon" If so, it sets query to return all pokemon in the database
	if r.URL.Path == "/pokemon" || id == "" {
		query = "SELECT number, name FROM pokemon;"
	} else {
		query = fmt.Sprintf("SELECT number, name FROM pokemon WHERE number=%s", id)
	}

	// Runs the query, and then defers closing the query until the end of the function.
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error Getting pokemon")
		log.Fatal(err)
	}
	defer rows.Close()

	// Extracts the data from the query, and places it into a slice of allPokemon
	for rows.Next() {
		var number int
		var name string

		err := rows.Scan(&number, &name)
		if err != nil {
			log.Fatal(err)
		}

		thisPokemon := Pokemon{
			Number: number,
			Name:   name,
		}

		allPokemon = append(allPokemon, thisPokemon)
	}
	return WriteJSON(w, http.StatusOK, allPokemon)
}
