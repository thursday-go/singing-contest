package models

import (
	"database/sql"
	"crypto/rand"
	"encoding/hex"

	_ "modernc.org/sqlite"
)

type Game struct {
	ID   string
	Name string
}

func LoadGameById(db *sql.DB, id string) (*Game, error) {
	var g Game
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

func LoadGames(db *sql.DB) ([]Game, error) {
	q := "SELECT id, name FROM games;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gs []Game
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func CreateGame(db *sql.DB, name string) (Game, error) {
	g := Game{
		ID:   randomId(3),
		Name: name,
	}
	q := "INSERT INTO games (id, name) VALUES (?, ?);"
	if _, err := db.Exec(q, g.ID, g.Name); err != nil {
		return Game{}, err
	}
	return g, nil
}

func DeleteGameByID(db *sql.DB, id string) error {
	q := "DELETE FROM games WHERE id = ?;"
	if _, err := db.Exec(q, id); err != nil {
		return err
	}
	return nil
}


func InitDB(db *sql.DB) error {
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
