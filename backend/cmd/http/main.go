package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/infra/db"
)

func main() {
	ctx := context.Background()
	env, ok := os.LookupEnv("POSTGRES_DBURL")
	if !ok {
		log.Fatal("POSTGRES_DBURL tidak ada")
	}
	fmt.Println(env)

	pool, err := db.OpenPgxPool(ctx,db.PgxConfig{
		DbURL:             os.Getenv("POSTGRES_DBURL"),
		MaxConns:        20,
		MinConns:        5,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
		HealthTimeout:   3 * time.Second,
	})

	if err != nil {
		log.Fatalf("DB init failed: %v",err)
		log.Fatal("Exiting....")
	}
	log.Println("DB init success")
	defer pool.Close()
}