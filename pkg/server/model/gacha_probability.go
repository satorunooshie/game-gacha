package model

import (
	"database/sql"
)

type GachaProbability struct {
	CollectionItemID string
	Ratio            int
}
type gachaProbabilityRepository struct {
	Conn *sql.DB
}
type GachaProbabilityRepositoryInterface interface {
	SelectGachaProbabilities() ([]*GachaProbability, error)
}

var _ GachaProbabilityRepositoryInterface = (*gachaProbabilityRepository)(nil)

func NewGachaProbabilityRepository(conn *sql.DB) *gachaProbabilityRepository {
	return &gachaProbabilityRepository{
		Conn: conn,
	}
}

func (r *gachaProbabilityRepository) SelectGachaProbabilities() ([]*GachaProbability, error) {
	rows, err := r.Conn.Query("SELECT * FROM gacha_probabilities")
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
