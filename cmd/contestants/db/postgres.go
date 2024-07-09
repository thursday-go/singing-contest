package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

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
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")

	createTable()
	insertSampleData()
}

func createTable() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS contestants (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		location TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertSampleData() {
	_, err := DB.Exec(`INSERT INTO contestants (name, location) VALUES
		('Alice', 'Studio City'), ('Bob', 'Los Angeles'), ('Charlie', 'La Canada')
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatal(err)
	}
}
