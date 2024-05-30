package model

import "time"

type Item struct {
	ID uint `json:"id" gorm:"primaryKey"`

	CreatedAt time.Time

	Title string `json:"item_title"`

	Completed bool `json:"completed"`

	ListID uint `json:"list_id"`
}
