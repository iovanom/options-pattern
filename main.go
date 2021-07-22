package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ServerOption func(*http.Server)

func NewServer(opts ...ServerOption) *http.Server {
	const (
		defaultAddr         = ":8080"
		defaultReadTimeout  = 5 * time.Second
		defaultWriteTimeout = 5 * time.Second
	)

	srv := &http.Server{
		Addr:         defaultAddr,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func WithTimeouts(readTimeout time.Duration, writeTimeout time.Duration) ServerOption {
	return func(s *http.Server) {
		WithReadTimeout(readTimeout)(s)
		WithWriteTimeout(writeTimeout)(s)
	}
}

func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(s *http.Server) {
		s.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(s *http.Server) {
		s.WriteTimeout = timeout
	}
}

func WithAddress(address string) ServerOption {
	return func(s *http.Server) {
		s.Addr = address
	}
}

func WithHandler(handler http.Handler) ServerOption {
	return func(s *http.Server) {
		s.Handler = handler
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})
	srv := NewServer(
		WithAddress(":8099"),
		WithHandler(mux),
		WithTimeouts(15*time.Second, 15*time.Second),
	)
	log.Fatal(srv.ListenAndServe())
}
