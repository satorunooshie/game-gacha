package model

import (
	"database/sql"
	"fmt"
	"strings"
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
func SaveUserCollectionItems(tx *sql.Tx, newItemIDs []string, userID string) error {
	placeholder := strings.Repeat("(?, ?, ?),", len(newItemIDs))
	queryArgs := make([]interface{}, 0, len(newItemIDs)*3)
	for _, itemID := range newItemIDs {
		queryArgs = append(queryArgs, userID, itemID, time.Now())
	}
	joinedQuery := strings.Trim(fmt.Sprintf("INSERT INTO user_collection_items(user_id, collection_item_id, created_at) VALUES %s", placeholder), ",")
	stmt, err := tx.Prepare(joinedQuery)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(queryArgs...); err != nil {
		return err
	}
	return nil
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
