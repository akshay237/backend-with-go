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
	"sync"
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
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/reflection"
)

const (
	HTTP    = "HTTP"
	GRPC    = "GRPC"
	GATEWAY = "GATEWAY"
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
	ready := make(chan struct{})
	var httpSrv *http.Server
	var grpcSrv *grpc.Server

	go func() {
		log.Println("Application server is starting")
		var err error
		switch config.ServerName {
		case HTTP:
			httpSrv, err = runGinServer(config, store)
		case GRPC:
			grpcSrv, err = runGRPCServer(config, store)
		case GATEWAY:
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				httpSrv, _ = runGatewayServer(config, store)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				grpcSrv, err = runGRPCServer(config, store)
			}()
			wg.Wait()
			if err != nil {
				log.Fatalf("failed to start server: %v", err)
			}
		default:
			err = fmt.Errorf("unknown server name: %s", config.ServerName)
			log.Fatalf("failed to start server: %v", err)
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
	case <-time.After(30 * time.Second):
		log.Println("Server startup timed out")
	}

	switch config.ServerName {
	case HTTP:
		if httpSrv != nil {
			log.Println("Stopping HTTP server")
			err := httpSrv.Shutdown(context.Background())
			if err != nil {
				log.Println("HTTP Server Shutdown failed:", err)
			}
		}
	case GRPC:
		if grpcSrv != nil {
			log.Println("Stopping gRPC server")
			grpcSrv.GracefulStop()
		}
	case GATEWAY:
		if httpSrv != nil {
			log.Println("Stopping gRPC Gateway server")
			err := httpSrv.Shutdown(context.Background())
			if err != nil {
				log.Println("gRPC Gateway Server Shutdown failed:", err)
			}
		}
		if grpcSrv != nil {
			log.Println("Stopping gRPC server")
			grpcSrv.GracefulStop()
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

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	return srv, nil
}

func runGRPCServer(config util.Config, store db.Store) (*grpc.Server, error) {
	handler, err := gapi.NewServerHandler(config, store)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server handler: %v", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, handler)
	reflection.Register(gRPCServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %v", config.GRPCServerAddress, err)
	}

	log.Println("Starting gRPC server on", config.GRPCServerAddress)
	err = gRPCServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}

	return gRPCServer, nil
}

func runGatewayServer(config util.Config, store db.Store) (*http.Server, error) {
	handler, err := gapi.NewServerHandler(config, store)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server handler: %v", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to register gRPC gateway: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	srv := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: mux,
	}

	log.Println("Starting HTTP Gateway server on", config.HTTPServerAddress)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	// lis, err := net.Listen("tcp", config.HTTPServerAddress)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to listen on %s: %v", config.HTTPServerAddress, err)
	// }

	// log.Println("Starting gRPC server on", config.GRPCServerAddress)
	// err = http.Serve(lis, mux)
	// if err != nil {
	// 	log.Fatalf("failed to serve gRPC server: %v", err)
	// }

	return srv, nil
}
