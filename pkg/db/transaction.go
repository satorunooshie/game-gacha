package db

import (
	"database/sql"
	"log"
)

type database struct {
	DB *sql.DB
}

var (
	DB database
)

func (db *database) Transaction(txFunc func(*sql.Tx) error) (err error) {
	tx, err := Conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if derr := tx.Rollback(); derr != nil {
				log.Println(derr, "rollback")
			}
			log.Panic(p)
		}
		if err != nil {
			if derr := tx.Rollback(); derr != nil {
				log.Println(derr, "rollback")
			}
			return
		}
		if derr := tx.Commit(); derr != nil {
			log.Println(derr, "commit")
		}
		return
	}()
	err = txFunc(tx)
	return
}
