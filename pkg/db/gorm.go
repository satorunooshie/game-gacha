package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"game-gacha/pkg/env"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Gorm *gorm.DB

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

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		database)
	if Gorm, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}); err != nil {
		log.Fatal("gorm connection failed", err)
	}
}
