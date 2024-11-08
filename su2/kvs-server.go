package su2

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Storage map[string]string

func (s *Storage) Set(key string, value string) {
	(*s)[key] = value
}

func (s *Storage) Get(key string) string {
	return (*s)[key]
}

var (
	address string
	port    string
	storage = Storage{}
)

func KvsServer() {
	flag.StringVar(&port, "p", "8080", "port to listen on")
	flag.StringVar(&address, "a", "localhost", "address to listen on")

	flag.Parse()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              net.JoinHostPort(address, port),
		Handler:           mux,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	RunKvsServer(mux, server)
}

func RunKvsServer(mux *http.ServeMux, s *http.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, os.Kill)

	mux.HandleFunc("GET /get", GetValueHandler)
	mux.HandleFunc("POST /set", SetValueHandler)

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

func SetValueHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := r.FormValue("key")
	value := r.FormValue("value")

	if len(key) == 0 || len(value) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isValueExist := len(storage.Get(key)) != 0

	storage.Set(key, value)

	if isValueExist {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetValueHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	key := params.Get("key")

	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := storage.Get(key)

	if len(value) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
