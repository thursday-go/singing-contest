package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"singing-contest/models"
)

const (
	dbPath  = ":memory:"
	svrPath = ":8080"
)

func main() {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Add some fake data
	if err := models.InitDB(db); err != nil {
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
		gs, err := models.LoadGames(db)
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
		g, err := models.CreateGame(db, gname)
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
		g, err := models.LoadGameById(db, gid)
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
		if err := models.DeleteGameByID(db, gid); err != nil {
			panic(err)
		}

		// Return empty response
		w.WriteHeader(http.StatusOK)
	})

	// Run the server
	fmt.Println("Starting server", svrPath)
	if err := http.ListenAndServe(svrPath, mux); err != nil {
		panic(err)
	}
}
