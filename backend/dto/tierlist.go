package dto

import "time"

type CreateTierlistRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	ExpiryTime  time.Time           `json:"expiry_time" binding:"required"`
	Items       []CreateItemRequest `json:"tierlist_items"`
}

type CreateItemRequest struct {
	Name      string `json:"name" binding:"required"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

type UpdateItemRequest struct {
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}
