package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"contestants-service/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")

	createTable()
	insertSampleData()
}

func CloseDB() {
	db.Close()
}

func FetchContestantsFromDB() ([]models.Contestant, error) {
	rows, err := db.Query("SELECT id, name, location FROM contestants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contestants []models.Contestant
	for rows.Next() {
		var c models.Contestant
		if err := rows.Scan(&c.ID, &c.Name, &c.Location); err != nil {
			return nil, err
		}
		contestants = append(contestants, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contestants, nil
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS contestants (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		location TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertSampleData() {
	_, err := db.Exec(`INSERT INTO contestants (name, location) VALUES
		('Alice', 'Studio City'), ('Bob', 'Los Angeles'), ('Charlie', 'La Canada')
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatal(err)
	}
}
