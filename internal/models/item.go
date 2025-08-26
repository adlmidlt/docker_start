package models

import (
	"database/sql"
	"time"
)

type Item struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func ScanItem(row *sql.Row) (*Item, error) {
	var item Item
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.CreateAt,
		&item.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func ScanItems(rows *sql.Rows) ([]*Item, error) {
	var items []*Item
	if rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.CreateAt,
			&item.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
