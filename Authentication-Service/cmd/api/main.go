package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var counts int64

const webPort = "80"

func main() {
	log.Println("Starting authentication service...")
	conn := connectToDb()
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	app.routes()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
	defer closeDB(conn)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")
	fmt.Println("Connecting to database...")
	fmt.Println(dsn)
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Databse not yet ready..")
			counts++
		} else {
			log.Println("Connected to database")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
