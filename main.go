package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/pokemon", handleGetAllPokemon) // Should only be doing GET requests on this route for all pokemon
	http.HandleFunc("/pokemon/", handlePokemon)      // Emulates "/pokemon/{id}" on a framework as this is a limitation of only using net/http
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
	// Opens the pokemon database & defers closing until the end of the function
	db, errs := sql.Open("sqlite", "./pokemon.db")
	if errs != nil {
		log.Fatal(errs)
	}
	defer db.Close()

	// Switch statement to route the request depending on the method used
	switch r.Method {
	case "GET":
		handleGetPokemonByNumber(w, r, db)
	case "POST":
		handlePostPokemon(w, r, db)
	case "PATCH":
		w.Write([]byte("This is a patch request"))
	case "DELETE":
		w.Write([]byte("This is a delete request"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sorry this method is not supported"))
	}
}

func handleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	// Opens the pokemon database & defers closing until the end of the function
	db, errs := sql.Open("sqlite", "./pokemon.db")
	if errs != nil {
		log.Fatal(errs)
	}
	defer db.Close()

	allPokemon := []Pokemon{}

	// Runs the query, and then defers closing the query until the end of the function.
	rows, err := db.Query("SELECT number, name FROM pokemon;")
	if err != nil {
		fmt.Println("Error getting pokemon")
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
	WriteJSON(w, http.StatusOK, allPokemon)
}

func handleGetPokemonByNumber(w http.ResponseWriter, r *http.Request, db *sql.DB) any {
	// Extracts the id from the URL as a string
	numString := strings.TrimPrefix(r.URL.Path, "/pokemon/")

	//convert idString to int id and returns bad request if it is not a number
	num, err := strconv.Atoi(numString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter a valid Number"))
		return nil
	}

	// Gets maximum number of pokemon in database
	var maxNumber int
	query := db.QueryRow("SELECT MAX(number) FROM pokemon;")
	err = query.Scan(&maxNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Checks if the id is valid and returns bad request if it is not
	if num < 1 || num > maxNumber {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter a valid Number (Between 1 and " + strconv.Itoa(maxNumber) + ")"))
		return nil
	}

	// Runs the query on a single row of the database.
	row := db.QueryRow("SELECT number, name FROM pokemon WHERE number=?", num)

	var number int
	var name string

	// Scans the row and places the data into the variable thisPokemon
	e := row.Scan(&number, &name)
	if e != nil {
		log.Fatal(e)
	}

	thisPokemon := Pokemon{
		Number: number,
		Name:   name,
	}
	return WriteJSON(w, http.StatusOK, thisPokemon)
}

func handlePostPokemon(w http.ResponseWriter, r *http.Request, db *sql.DB) any {
	// Extracts the data from the POST request
	var pokemon Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		log.Fatal(err)
	}

	// Inserts the new pokemon into the database
	query := fmt.Sprintf("INSERT INTO pokemon VALUES (%d, '%s')", pokemon.Number, pokemon.Name)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Returns http status 200
	return http.StatusOK
}
