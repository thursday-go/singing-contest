package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"

	_ "modernc.org/sqlite"
)

const (
	dbPath  = ":memory:"
	svrPath = ":8080"
)

func initDB(db *sql.DB) error {
	queries := []string{
		"CREATE TABLE games (id TEXT PRIMARY KEY, name TEXT);",
		"INSERT INTO games (id, name) VALUES ('abc123', 'First Game');",
		"INSERT INTO games (id, name) VALUES ('def456', 'Second Game');",
		"INSERT INTO games (id, name) VALUES ('ghi789', 'Third Game');",
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func randomId(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type game struct {
	ID   string
	Name string
}

func loadGameById(db *sql.DB, id string) (*game, error) {
	var g game
	q := "SELECT id, name FROM games WHERE id = ?;"
	err := db.QueryRow(q, id).Scan(&g.ID, &g.Name)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &g, nil
	}
}

func loadGames(db *sql.DB) ([]game, error) {
	q := "SELECT id, name FROM games;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gs []game
	for rows.Next() {
		var g game
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func createGame(db *sql.DB, name string) (game, error) {
	g := game{
		ID:   randomId(3),
		Name: name,
	}
	q := "INSERT INTO games (id, name) VALUES (?, ?);"
	if _, err := db.Exec(q, g.ID, g.Name); err != nil {
		return game{}, err
	}
	return g, nil
}

func deleteGameByID(db *sql.DB, id string) error {
	q := "DELETE FROM games WHERE id = ?;"
	if _, err := db.Exec(q, id); err != nil {
		return err
	}
	return nil
}

func main() {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Add some fake data
	if err := initDB(db); err != nil {
		panic(err)
	}

	// Load the templates
	tmpls, err := template.ParseGlob("./templates/**/*.tmpl")
	if err != nil {
		panic(err)
	}

	// Create a server
	mux := http.NewServeMux()

	// Add some handlers
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// Load the games
		gs, err := loadGames(db)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "text/html")
		tmpls.ExecuteTemplate(w, "index", gs)
	})
	mux.HandleFunc("POST /games/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		// Get the gameid param from the URL
		gname := r.FormValue("name")

		// Create it in the DB
		g, err := createGame(db, gname)
		if err != nil {
			panic(err)
		}

		// Return a partial for HTMX
		w.Header().Set("Content-Type", "text/html")
		tmpls.ExecuteTemplate(w, "new-game", g)
	})
	mux.HandleFunc("GET /games/{gameid}/{$}", func(w http.ResponseWriter, r *http.Request) {
		// Get the gameid param from the URL
		gid := r.PathValue("gameid")

		// Load the game from the db
		g, err := loadGameById(db, gid)
		if err != nil {
			panic(err)
		}

		// Not found?
		if g == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Render the page
		w.Header().Set("Content-Type", "text/html")
		tmpls.ExecuteTemplate(w, "game", g)
	})
	mux.HandleFunc("DELETE /games/{gameid}/{$}", func(w http.ResponseWriter, r *http.Request) {
		// Get the gameid param from the URL
		gid := r.PathValue("gameid")

		// Delete the game from the db
		if err := deleteGameByID(db, gid); err != nil {
			panic(err)
		}

		// Return empty response
		w.WriteHeader(http.StatusOK)
	})

	// Run the server
	if err := http.ListenAndServe(svrPath, mux); err != nil {
		panic(err)
	}
}
