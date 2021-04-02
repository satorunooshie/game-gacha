package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/env"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

func init() {
	var err error
	err = env.Load()
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	// user:password@tcp(host:port)/database
	Conn, err = sql.Open(constant.DriverName,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", user, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}
}
