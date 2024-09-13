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
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	if err := godotenv.Load("cmd/main.env"); err != nil {
		log.Fatal("No .env file found")
	}
	SERVER_ADDRESS := os.Getenv("SERVER_ADDRESS")
	POSTGRES_CONN := os.Getenv("POSTGRES_CONN")

	db, err := pgxpool.New(context.Background(), POSTGRES_CONN)
	if err != nil {
		fmt.Println("error wih db", err)
	}

	userRepo := repositories.NewUserRepo(db)
	tenderRepo := repositories.NewTenderRepo(db)
	tenderService := services.NewTenderService(tenderRepo, userRepo)
	tender := delivery.NewTenderHandler(tenderService, logger)
	bidRepo := repositories.NewBidRepo(db)
	bidService := services.NewBidService(bidRepo, userRepo, tenderRepo)
	bid := delivery.NewBidHandler(bidService, logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/tenders", tender.GetTenderList).Methods("GET")
	r.HandleFunc("/tenders/new", tender.CreateTender).Methods("POST")
	r.HandleFunc("/tenders/my", tender.GetTenderByUser).Methods("GET")
	r.HandleFunc("/tenders/{tenderId}/status", tender.GetTenderStatus).Methods("GET")
	r.HandleFunc("/tenders/{tenderId}/status", tender.ChangeTenderStatus).Methods("PUT")
	r.HandleFunc("/tenders/{tenderId}/edit", tender.EditTender).Methods("PATCH")
	r.HandleFunc("/bids/new", bid.CreateBid).Methods("POST")
	r.HandleFunc("/bids/my", bid.GetUserBids).Methods("GET")
	r.HandleFunc("/bids/{tenderId}/my", bid.GetBidsForTender).Methods("GET")
	r.HandleFunc("/bids/{bidId}/status", bid.GetBidStatus).Methods("GET")
	r.HandleFunc("/bids/{bidId}/status", bid.ChangeBidStatus).Methods("PUT")
	r.HandleFunc("/bids/{bidId}/submit_decision", bid.SubmitBid).Methods("PUT")
	r.HandleFunc("/bids/{bidId}/edit", bid.EditBid).Methods("PATCH")

	srv := &http.Server{
		Addr:    SERVER_ADDRESS,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
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
