package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var serviceAvailable = true

type RestServer struct {
	Router *mux.Router
}

func NewRestServer() *RestServer {
	rs := &RestServer{}
	rs.Router = mux.NewRouter()
	Check(rs.Router)
	return rs
}

func ListenAndServe() {
	serviceName := "GitHook Relay"
	port := 8080
	address := fmt.Sprintf(":%d", port)
	rs := NewRestServer()
	server := &http.Server{
		Addr:         address,
		Handler:      rs.Router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	GracefulShutdown(server, serviceName)
	log.Printf("[INFO] [%s] Server starting on port %d", serviceName, port)
	err := server.ListenAndServe()
	if nil != err {
		log.Printf("[INFO] [%s] %s", serviceName, err.Error())
	}
}

func Check(router *mux.Router) {
	check := func(w http.ResponseWriter, req *http.Request) {
		if !serviceAvailable {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	metrics := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	router.HandleFunc("/check", check).Methods("GET")
	router.HandleFunc("/metrics", metrics).Methods("GET")
}

func GracefulShutdown(server *http.Server, serviceName string) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		signal.Stop(c)

		serviceAvailable = false

		log.Printf("[INFO] [%s] Server shutting down gracefully", serviceName)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		err := server.Shutdown(ctx)
		if nil != err {
			log.Printf("[ERROR] [%s] Error shutting down server: %s", serviceName, err.Error())
		}
	}()
}
