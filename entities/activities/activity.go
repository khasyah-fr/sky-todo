package activities

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ActivityID int    `gorm:"primaryKey;autoIncrement"`
	Title      string `gorm:"not null"`
	Email      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
