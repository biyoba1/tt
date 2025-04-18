package initializer

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

func PingDatabase() {
	s := os.Getenv("PG_DSN")
	if s == "" {
		panic("PG_DSN is not set in the environment variables")
	}
	db, err := sqlx.Open("postgres", s)
	if err != nil {
		panic("failed open database")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("bombardiro crocodilo")
}
