package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Db *sqlx.DB

func init() {
	connectionParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("RDB_HOST"),
		os.Getenv("RDB_PORT"),
		os.Getenv("RDB_USERNAME"),
		os.Getenv("RDB_PASSWORD"),
		os.Getenv("RDB_DBNAME"),
	)

	db, err := sqlx.Connect("postgres", connectionParams)
	if err != nil {
		log.Fatal(err)
	}
	Db = db
}

func SaveIpToDatabase(ip string) {
	if ip != "" {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		Db.MustExec(
			"INSERT INTO views (client_ip, view_date) VALUES ($1, $2)", ip, currentTime,
		)
	}
}
