package model

import (
	"database/sql"
	"time"

	"game-gacha/pkg/db"
)

type CollectionItem struct {
	ID        string
	Name      string
	Rarity    int
	CreatedAt *time.Time
}

func SelectAllCollectionItems() ([]*CollectionItem, error) {
	rows, err := db.Conn.Query("SELECT * FROM collection_items")
	if err != nil {
		return nil, err
	}
	return convertToCollectionItems(rows)
}
func convertToCollectionItems(rows *sql.Rows) ([]*CollectionItem, error) {
	defer rows.Close()
	collectionItems := make([]*CollectionItem, 0)
	for rows.Next() {
		var collectionItem CollectionItem
		if err := rows.Scan(&collectionItem.ID, &collectionItem.Name, &collectionItem.Rarity, &collectionItem.CreatedAt); err != nil {
			return nil, err
		}
		collectionItems = append(collectionItems, &collectionItem)
	}
	return collectionItems, nil
}
