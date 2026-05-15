package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tierlist struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreatorID uuid.UUID `json:"creator_id" gorm:"type:uuid;index"`
	ShareCode string `json:"share_code" gorm:"uniqueIndex"`
	ExpiryTime time.Time `json:"expiry_time"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // soft delete field
	Creator User `json:"creator" gorm:"foreignKey:CreatorID"` // association to the User model
	TierlistItems []TierlistItem `json:"tierlist_items" gorm:"foreignKey:TierlistID;constraint:OnDelete:CASCADE"` // association to the TierlistItem model
	Submissions []Submissions `json:"submissions" gorm:"foreignKey:TierlistID;constraint:OnDelete:CASCADE"` // association to the Submissions model
}

type TierlistItem struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TierlistID uuid.UUID `json:"tierlist_id" gorm:"type:uuid;index"`
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
	SortOrder int `json:"sort_order"`
}

type Submissions struct { // this acts as the metadata table
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TierlistID uuid.UUID `json:"tierlist_id" gorm:"type:uuid;uniqueIndex:uq_tierlist_user;"` // composite unique index with UserID to ensure one submission per user per tierlist
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex:uq_tierlist_user"` // composite unique index with TierlistID to ensure one submission per user per tierlist
	CreatedAt time.Time `json:"created_at"`
	Rankings []SubmissionRankings `json:"rankings" gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE"` // association to the SubmissionRankings model
}

type SubmissionRankings struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SubmissionID uuid.UUID `json:"submission_id" gorm:"type:uuid;index"`
	ItemID uuid.UUID `json:"item_id" gorm:"type:uuid;index"`
	Tier string `json:"tier" gorm:"type:varchar(1)"` // e.g., "S", "A", "B", etc.
}