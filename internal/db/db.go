package db

import (
	"fmt"
	"log"
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

func UpdateBlackList(ip string) {
	Db.MustExec(
		"INSERT INTO blacklist (ip, day) VALUES ($1, $2)", ip, time.Now().Format("2006-01-02"),
	)
}

func GetIpsForBlackList() []string {
	today := time.Now()
	todayString := today.Format("2006-01-02")
	tomorrow := today.Add(time.Hour * 24)
	tomorrowString := tomorrow.Format("2006-01-02")

	ips := make([]string, 0)
	Db.Select(
		&ips,
		`
		SELECT client_ip FROM views
			WHERE view_date >= $1 AND view_date <= $2
				AND client_ip NOT IN ( SELECT ip FROM blacklist WHERE day = $1 )
			GROUP BY client_ip
			HAVING COUNT(1) >= 10
		`,
		todayString, tomorrowString,
	)
	return ips
}
