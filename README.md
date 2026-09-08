# Pokedex

A simple REST API for managing a Pokemon collection, written in Go using the standard `net/http` package and SQLite for storage.

## Tech Stack

- Go 1.18
- [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) — pure Go SQLite driver
- [`github.com/monzo/terrors`](https://github.com/monzo/terrors) — structured error handling

## Getting Started

### Prerequisites

- Go 1.18 or later

### Running the server

```bash
go run main.go
```

The server listens on `http://localhost:8080`. It reads from and writes to `test-pokemon.db` in the project root using the schema in `pokemon.sql`.

## API

Each `Pokemon` object has the shape:

```json
{
  "id": 1,
  "number": 1,
  "name": "Bulbasaur",
  "sprite": "https://..."
}
```

| Method | Route          | Description                                    |
| ------ | -------------- | ------------------------------------------------ |
| GET    | `/pokemon`     | Get all pokemon                                |
| GET    | `/pokemon/{id}`| Get a single pokemon by ID                     |
| POST   | `/pokemon/`    | Create a new pokemon (ID is supplied in the JSON body) |
| PATCH  | `/pokemon/{id}`| Update an existing pokemon                     |
| DELETE | `/pokemon/{id}`| Delete a pokemon                               |

## Database

The `pokemon` table is defined in `pokemon.sql`:

```sql
CREATE TABLE pokemon (
  id INTEGER PRIMARY KEY,
  number INTEGER NOT NULL,
  name TEXT NOT NULL,
  sprite VARCHAR(300) NOT NULL
);
```

`pokemon.sql` also contains commented-out seed data for the original 151 Pokemon, using sprite artwork from [PokeAPI/sprites](https://github.com/PokeAPI/sprites).

## Notes

- CORS is currently configured to allow requests only from a hardcoded local IP (`enableCors` in `main.go`) — update this to match your frontend's origin if running locally.
