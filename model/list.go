package model

import (
	"time"
)

type List struct {
	ID uint `json:"id" gorm:"primaryKey"`

	CreatedAt time.Time

	UserID uint `json:"user_id"`

	User User

	Title string `json:"list_title"`

	Progress float32 `json:"progress"`

	Items []Item `json:"list_items"`
}
