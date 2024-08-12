package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/juliofilizzola/server/internal/api"
	"github.com/juliofilizzola/server/internal/store/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	ctx := context.Background()
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		user, password, host, port, dbName,
	)
	pool, err := pgxpool.New(ctx, connString)

	if err != nil {
		panic(err)
	}

	defer func(p *pgxpool.Pool) {
		p.Close()
	}(pool)

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrAbortHandler) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
}
