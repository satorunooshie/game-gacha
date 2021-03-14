package model

import (
	"database/sql"

	"game-gacha/pkg/db"
)

type GachaProbability struct {
	CollectionItemID string
	Ratio            int
}

func SelectGachaProbabilities() ([]*GachaProbability, error) {
	rows, err := db.Conn.Query("SELECT * FROM gacha_probabilities")
	if err != nil {
		return nil, err
	}
	return convertToGachaProbabilities(rows)
}
func convertToGachaProbabilities(rows *sql.Rows) ([]*GachaProbability, error) {
	defer rows.Close()
	gachaProbabilities := make([]*GachaProbability, 0)
	for rows.Next() {
		var gachaProbability GachaProbability
		if err := rows.Scan(&gachaProbability.CollectionItemID, &gachaProbability.Ratio); err != nil {
			return nil, err
		}
		gachaProbabilities = append(gachaProbabilities, &gachaProbability)
	}
	return gachaProbabilities, nil
}
