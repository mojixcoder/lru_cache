package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MojtabaArezoomand/lru_cache/internal/config"
	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
)

// newRouter initializes a new router.
func newRouter() *mux.Router {
	r := mux.NewRouter()

	app := newApp()

	r.HandleFunc("/get/{key}", app.Get).Methods(http.MethodGet)
	r.HandleFunc("/set", app.Set).Methods(http.MethodPost)

	return r
}

// RunServer runs the server
func RunServer() {
	var cfg config.ServerConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	r := newRouter()

	srv := &http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Couldn't shutdown the server:", err)
	}
}
