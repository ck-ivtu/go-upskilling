package su2

import (
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ChiKvsServer() {
	flag.StringVar(&port, "p", "8080", "port to listen on")
	flag.StringVar(&address, "a", "localhost", "address to listen on")

	flag.Parse()

	router := chi.NewRouter()

	server := &http.Server{
		Addr:              net.JoinHostPort(address, port),
		Handler:           router,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	RunChiKvsServer(router, server)
}

func RunChiKvsServer(router chi.Router, s *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	router.Get("/get", GetValueHandler)
	router.Post("/set", SetValueHandler)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			return
		}
	}()

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s\n", err.Error())
	}

	log.Printf("Server exiting\n")
}
