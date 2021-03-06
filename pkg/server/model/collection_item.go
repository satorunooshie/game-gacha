package model

import (
	"database/sql"
	"time"
)

type CollectionItem struct {
	ID        string
	Name      string
	Rarity    int
	CreatedAt *time.Time
}
type collectionItemRepository struct {
	Conn *sql.DB
}
type CollectionItemRepositoryInterface interface {
	SelectAllCollectionItems() ([]*CollectionItem, error)
}

var _ CollectionItemRepositoryInterface = (*collectionItemRepository)(nil)

func NewCollectionItemRepository(conn *sql.DB) *collectionItemRepository {
	return &collectionItemRepository{
		Conn: conn,
	}
}

func (r *collectionItemRepository) SelectAllCollectionItems() ([]*CollectionItem, error) {
	rows, err := r.Conn.Query("SELECT * FROM collection_items")
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
