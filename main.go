package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/akshay237/backend-with-go/api"
	db "github.com/akshay237/backend-with-go/database/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	// 0. Create a database connection and pass it to the server
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("db connection failed: ", err)
	}

	// 1. create the database store using the db connection
	store := db.NewStore(conn)

	// 2. create an server and start the server
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("failed to start the server", err)
	}
}
