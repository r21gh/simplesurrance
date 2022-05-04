package main

import (
	"context"
	"flag"
	"github.com/r21gh/simplesurrance/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	serverAddress string
)

func main() {
	// parse flags and set server address
	flag.StringVar(&serverAddress, services.ServerName, services.ServerPortWithColon, services.ServerUsage)
	flag.Parse()

	// create a new logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println(services.ServerWelcomeMessage)

	// create a new router using the standard library
	router := http.NewServeMux()

	// add the counter handler
	router.Handle(services.CounterPath, services.ApiHandler())

	// define a new server configuration using the standard library
	server := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       15 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Addr:              serverAddress,
		Handler:           services.Tracing(services.NewRequestID)(services.Logging(logger)(router)),
		ErrorLog:          logger,
	}

	// `done` is a channel that is closed when the server is shutdown
	done := make(chan bool)

	// quit channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// shutdown the server when the quit channel is closed
	go func() {
		<-quit
		logger.Println(services.ServerShutdownMessage)

		// create a context with a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// shutdown the server gracefully
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("error in shutting down the server gracefully %v", err)
		}
		close(done)
	}()

	logger.Println("Server is working on", serverAddress)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", serverAddress, err)
	}

	<-done
	logger.Println("Server stopped")
}
