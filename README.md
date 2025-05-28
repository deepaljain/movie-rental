# Movie Rental API

A simple RESTful API for browsing movies and managing a rental cart, built with Go, Gin, and PostgreSQL.

## Features

- List all movies, with filters for genre, actor, and year
- Get movie details by ID
- Add movies to a user's cart
- View a user's cart
- Simple hello endpoint for testing

## Project Structure

```
movie-rental/
├── cmd/                # Main application entrypoint
├── pkg/
│   ├── movies/         # Movie handlers, models, tests
│   ├── cart/           # Cart handlers, models, tests
│   └── hello/          # Hello handler
├── migrations/         # Database migration SQL files
├── go.mod              # Go dependencies
├── Makefile            # Build/test/migrate commands
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- [migrate](https://github.com/golang-migrate/migrate) CLI tool

### Setup

1. **Clone the repository**
   ```sh
   git clone <your-repo-url>
   cd movie-rental
   ```

2. **Configure the database**
   - Create a PostgreSQL database named `movierental`.
   - Update the connection string in `cmd/main.go` if needed.

3. **Run migrations**
   ```sh
   make migrateup
   ```

4. **Insert sample data**
   - Use the provided SQL script to populate the `movies` table.

5. **Build and run the server**
   ```sh
   make build
   make run
   ```
   The API will be available at `http://localhost:8080`.

## API Endpoints

- `GET /hello` — Returns a hello message
- `GET /movies` — List all movies (supports `genre`, `actor`, `year` query params)
- `GET /movies/:id` — Get movie by ID
- `POST /cart` — Add a movie to a user's cart (JSON: `{ "user_id": int, "movie_id": int }`)
- `GET /cart/:user_id` — View a user's cart

## Example Usage

```sh
curl http://localhost:8080/movies
curl "http://localhost:8080/movies?genre=Action"
curl http://localhost:8080/movies/1
curl -X POST -H "Content-Type: application/json" -d '{"user_id":1,"movie_id":2}' http://localhost:8080/cart
curl http://localhost:8080/cart/1
```

## Running Tests

```sh
make test
```

## License

MIT