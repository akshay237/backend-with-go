package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/akshay237/backend-with-go/api"
	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/util"
)

func main() {

	// 0. load the config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load the config: ", err)
	}

	// 1. Create a database connection and pass it to the server
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("db connection failed: ", err)
	}

	// 2. create the database store using the db connection
	store := db.NewStore(conn)

	// 3. create an server and start the server
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("failed to start the server", err)
	}
}
