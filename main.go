package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/akshay237/backend-with-go/api"
	db "github.com/akshay237/backend-with-go/db/sqlc"
	"github.com/akshay237/backend-with-go/util"
)

func waitTillStopFile(stoppedflag *uint32, stopch chan string, stopfilepath string) {
	log.Println("Stopfile Path:", stopfilepath)
	for atomic.LoadUint32(stoppedflag) == 0 {
		if _, err := os.Stat(stopfilepath); err == nil {
			log.Println("Removing stop file")
			err := os.Remove(stopfilepath)
			if err != nil {
				log.Println("Removing stopfile has failed")
			}
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	stopch <- "stopfilereceived"
}

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
	errs := make(chan error)
	handler := api.NewServerHandler(store)
	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: handler.Router,
	}
	go func() {
		log.Println("Application server is starting")
		errs <- srv.ListenAndServe()
	}()

	// 4. Graceful shutdown of the application using os signal and stopfile
	go func() {
		ossignalch := make(chan os.Signal, 1)
		signal.Notify(ossignalch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		signalrcvd := <-ossignalch
		log.Println("Signal Receieved from OS:", signalrcvd.String())
		errs <- fmt.Errorf("%s", signalrcvd)
	}()

	stopfilech := make(chan string, 1)
	stoppedflag := uint32(0)

	go func() {
		waitTillStopFile(&stoppedflag, stopfilech, path.Join(config.StopFilePath, "stopfile"))
		atomic.StoreUint32(&stoppedflag, 1)
	}()

	select {
	case <-errs:
		log.Println("Recieved message on error chan")
	case <-stopfilech:
		log.Println("Stop file seen, so stop the server")
	}
	atomic.StoreUint32(&stoppedflag, 1)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server is shutdown successfully")
}
