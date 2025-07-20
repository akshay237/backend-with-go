package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/akshay237/backend-with-go/api"
	db "github.com/akshay237/backend-with-go/database/sqlc"
	"github.com/akshay237/backend-with-go/gapi"
	"github.com/akshay237/backend-with-go/pb"
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

	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 Recovered from panic in main: %v", r)
			os.Exit(1) // Ensure non-zero exit code so make knows it failed
		}
	}()

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
	ready := make(chan struct{})
	var httpSrv *http.Server
	var grpcSrv *grpc.Server

	go func() {
		log.Println("Application server is starting")
		var err error
		if config.RunGRPC {
			grpcSrv, err = runGRPCServer(config, store)
		} else {
			httpSrv, err = runGinServer(config, store)
		}
		close(ready)
		if err != nil {
			errs <- err
		}
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
	case err = <-errs:
		log.Println("Recieved message on error chan", err)
	case <-stopfilech:
		log.Println("Stop file seen, so stop the server")
	}
	atomic.StoreUint32(&stoppedflag, 1)

	select {
	case <-ready:
		log.Println("Server is ready to stop")
	case <-time.After(10 * time.Second):
		log.Println("Server startup timed out")
	}

	if config.RunGRPC && grpcSrv != nil {
		log.Println("Stopping gRPC server")
		grpcSrv.GracefulStop()
	} else if !config.RunGRPC && httpSrv != nil {
		log.Println("Stopping HTTP server")
		err := httpSrv.Shutdown(context.Background())
		if err != nil {
			log.Println("Server Shutdown:", err)
		}
	}
	log.Println("Server is shutdown successfully")
}

func runGinServer(config util.Config, store db.Store) (*http.Server, error) {
	handler, err := api.NewServerHandler(config, store)
	if err != nil {
		return nil, fmt.Errorf("cannot create server handler: %v", err)
	}

	srv := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: handler.Router,
	}

	log.Println("Starting Gin server on", config.HTTPServerAddress)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv, nil
}

func runGRPCServer(config util.Config, store db.Store) (*grpc.Server, error) {
	handler, err := gapi.NewServerHandler(config, store)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server handler: %v", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, handler)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %v", config.GRPCServerAddress, err)
	}

	log.Println("Starting gRPC server on", config.GRPCServerAddress)
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
	return gRPCServer, nil
}
