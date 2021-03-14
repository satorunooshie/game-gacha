package model

import (
	"database/sql"
	"time"

	"game-gacha/pkg/db"
)

type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
	CreatedAt        *time.Time
}

func SelectUserCollectionItems(userID string) ([]*UserCollectionItem, error) {
	rows, err := db.Conn.Query("SELECT * FROM user_collection_items WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	return convertToUserCollectionItems(rows)
}
func convertToUserCollectionItems(rows *sql.Rows) ([]*UserCollectionItem, error) {
	defer rows.Close()
	userCollectionItems := make([]*UserCollectionItem, 0)
	for rows.Next() {
		var userCollectionItem UserCollectionItem
		if err := rows.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionItemID, &userCollectionItem.CreatedAt); err != nil {
			return nil, err
		}
		userCollectionItems = append(userCollectionItems, &userCollectionItem)
	}
	return userCollectionItems, nil
}
