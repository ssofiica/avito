package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zadanie-6105/internal/delivery"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://svalova:mydbpass@localhost:5432/test-gaz?sslmode=disable")
	if err != nil {
		fmt.Println("error wih db", err)
	}
	fmt.Println("ok")

	tenderRepo := repositories.NewTenderRepo(db)
	tenderService := services.NewTenderService(tenderRepo)
	tender := delivery.NewTenderHandler(tenderService)
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/tenders", tender.GetTenderList).Methods("GET")
	r.HandleFunc("/tenders/new", nil).Methods("POST")
	r.HandleFunc("/tenders/my", nil).Methods("GET")
	r.HandleFunc("/tenders/{tenderId}/edit", nil).Methods("PATCH")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
