package model

import (
	"database/sql"
	"time"

	"game-gacha/pkg/db"
)

type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int
	Coin      int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func InsertUser(user *User) error {
	stmt, err := db.Conn.Prepare("INSERT INTO users(id, auth_token, name, high_score, coin, created_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.ID, user.AuthToken, user.Name, user.HighScore, user.Coin, time.Now()); err != nil {
		return err
	}
	return nil
}
func SelectUserByPK(userID string) (*User, error) {
	row := db.Conn.QueryRow("SELECT * FROM users WHERE id = ?", userID)
	return convertToUser(row)
}
func SelectUserByAuthToken(authToken string) (*User, error) {
	row := db.Conn.QueryRow("SELECT * FROM users WHERE auth_token = ?", authToken)
	return convertToUser(row)
}
func UpdateUserByPK(user *User) error {
	stmt, err := db.Conn.Prepare("UPDATE users SET name = ?, high_score = ?, coin = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(user.Name, user.HighScore, user.Coin, time.Now(), user.ID); err != nil {
		return err
	}
	return nil
}
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	if err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
