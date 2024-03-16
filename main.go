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
	Id     int    `json:"id"`
	Number int    `json:"number"`
	Name   string `json:"name"`
	Sprite string `json:"sprite"`
}

// Added log.Fatal so that it sends an error if the server crashes
func main() {
	http.HandleFunc("/pokemon", handleGetAllPokemon) // Should only be doing GET requests on this route for all pokemon
	http.HandleFunc("/pokemon/", handlePokemon)      // Emulates "/pokemon/{id}" on a framework as this is a limitation of only using net/http
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Encodes data into JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func enableCors(w *http.ResponseWriter) {
	// Enables CORS for the frontend to access the API
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
}

// Handles the incoming http requests and routes them depending on the method used.
func handlePokemon(w http.ResponseWriter, r *http.Request) {
	// Enables CORS
	enableCors(&w)

	// Opens the pokemon database & defers closing until the end of the function
	db, errs := sql.Open("sqlite", "./all-pokemon.db")
	if errs != nil {
		log.Fatal(errs)
	}
	defer db.Close()

	// Switch statement to route the request depending on the method used
	switch r.Method {
	case "GET":
		handleGetPokemonByNumber(w, r, db)
	case "POST":
		handlePostPokemon(&w, r, db)
	case "PATCH":
		w.Write([]byte("This is a patch request"))
	case "DELETE":
		w.Write([]byte("This is a delete request"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sorry this method is not supported"))
	}
}

// Handles the GET request for all pokemon
func handleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	// Enables CORS
	enableCors(&w)

	// Opens the pokemon database & defers closing until the end of the function
	db, errs := sql.Open("sqlite", "./all-pokemon.db")
	if errs != nil {
		log.Fatal(errs)
	}
	defer db.Close()

	allPokemon := []Pokemon{}

	// Runs the query, and then defers closing the query until the end of the function.
	rows, err := db.Query("SELECT * FROM pokemon;")
	if err != nil {
		fmt.Println("Error getting pokemon")
		log.Fatal(err)
	}
	defer rows.Close()

	// Extracts the data from the query, and places it into a slice of allPokemon
	for rows.Next() {

		thisPokemon := Pokemon{}

		err := rows.Scan(&thisPokemon.Id, &thisPokemon.Number, &thisPokemon.Name, &thisPokemon.Sprite)
		if err != nil {
			log.Fatal(err)
		}

		allPokemon = append(allPokemon, thisPokemon)
	}
	WriteJSON(w, http.StatusOK, allPokemon)
}

// Handles the GET request for a single pokemon
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
	row := db.QueryRow("SELECT * FROM pokemon WHERE number=?", num)

	thisPokemon := Pokemon{}

	// Scans the row and places the data into the variable thisPokemon
	e := row.Scan(&thisPokemon.Id, &thisPokemon.Number, &thisPokemon.Name, &thisPokemon.Sprite)
	if e != nil {
		log.Fatal(e)
	}

	return WriteJSON(w, http.StatusOK, thisPokemon)
}

// Handles the POST request for a single pokemon
func handlePostPokemon(w *http.ResponseWriter, r *http.Request, db *sql.DB) any {
	// Extracts the data from the POST request
	var pokemon Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		log.Fatal(err)
	}

	// Checks if the pokemon already exists in the database
	check := db.QueryRow("SELECT * FROM pokemon WHERE number=?", pokemon.Number)
	err = check.Scan(&pokemon.Id, &pokemon.Number, &pokemon.Name, &pokemon.Sprite)
	if err == nil {
		(*w).WriteHeader(http.StatusBadRequest)
		(*w).Write([]byte("Pokemon already exists"))
		return nil
	}

	// Adds the pokemon to the database
	query := fmt.Sprintf("INSERT INTO pokemon VALUES (%d, %d, '%s', '%s')", pokemon.Id, pokemon.Number, pokemon.Name, pokemon.Sprite)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Returns http status 200
	return http.StatusOK
}
