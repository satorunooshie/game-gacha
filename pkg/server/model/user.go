package model

import (
	"database/sql"
	"time"

	"game-gacha/pkg/constant"
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
type userRepository struct {
	Conn *sql.DB
}

// TODO: make interface smaller
type UserRepositoryInterface interface {
	InsertUser(user *User) error
	SelectUserByPK(userID string) (*User, error)
	SelectUserByPKForUpdate(tx *sql.Tx, userID string) (*User, error)
	SelectUserByAuthToken(authToken string) (*User, error)
	SelectUsersOrderByHighScore(startPosition, limit int) ([]*User, error)
	UpdateUserByPK(user *User) error
	UpdateUserCoinByPK(tx *sql.Tx, coin int, userID string) error
}

var _ UserRepositoryInterface = (*userRepository)(nil)

func NewUserRepository(conn *sql.DB) *userRepository {
	return &userRepository{
		Conn: conn,
	}
}

func (r *userRepository) InsertUser(user *User) error {
	stmt, err := r.Conn.Prepare("INSERT INTO users(id, auth_token, name, high_score, coin, created_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.ID, user.AuthToken, user.Name, user.HighScore, user.Coin, time.Now()); err != nil {
		return err
	}
	return nil
}
func (r *userRepository) SelectUserByPK(userID string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * FROM users WHERE id = ?", userID)
	return convertToUser(row)
}
func (r *userRepository) SelectUserByPKForUpdate(tx *sql.Tx, userID string) (*User, error) {
	row := tx.QueryRow("SELECT * FROM users WHERE id = ? FOR UPDATE", userID)
	return convertToUser(row)
}
func (r *userRepository) SelectUserByAuthToken(authToken string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * FROM users WHERE auth_token = ?", authToken)
	return convertToUser(row)
}
func (r *userRepository) SelectUsersOrderByHighScore(startPosition, limit int) ([]*User, error) {
	rows, err := r.Conn.Query("SELECT * FROM users WHERE high_score > 0 ORDER BY high_score DESC, id ASC LIMIT ? OFFSET ?", limit, startPosition)
	if err != nil {
		return nil, err
	}
	return convertToUsers(rows)
}
func (r *userRepository) UpdateUserByPK(user *User) error {
	stmt, err := r.Conn.Prepare("UPDATE users SET name = ?, high_score = ?, coin = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(user.Name, user.HighScore, user.Coin, time.Now(), user.ID); err != nil {
		return err
	}
	return nil
}
func (r *userRepository) UpdateUserCoinByPK(tx *sql.Tx, coin int, userID string) error {
	stmt, err := tx.Prepare("UPDATE users SET coin = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(coin, time.Now(), userID); err != nil {
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
func convertToUsers(rows *sql.Rows) ([]*User, error) {
	users := make([]*User, 0, constant.RankingLimit)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
