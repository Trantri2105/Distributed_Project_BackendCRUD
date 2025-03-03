package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func PostgresConnect() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DB_NAME"))
	log.Println("Connecting to PostgreSQL: ", psqlInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to database")
	return db
}
