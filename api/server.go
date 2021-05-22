package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/viper"
)

type Server struct {
	*http.Server
}

func NewServer() (*Server, error) {
	api, err := New()
	if err != nil {
		return nil, err
	}

	var addr string
	port := viper.GetString("port")

	// allow port to be set as localhost:3000 in env during development to avoid "accept incoming network connection" request on restarts
	if strings.Contains(port, ":") {
		addr = port
	} else {
		addr = ":" + port
	}

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server. Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server has been stopped")
}
