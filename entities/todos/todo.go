package todos

import (
	"time"
)

type Todo struct {
	TodoID          int    `gorm:"primaryKey;autoIncrement"`
	ActivityGroupID int    `gorm:"not null"`
	Title           string `gorm:"not null"`
	Priority        string `gorm:"not null"`
	IsActive        bool   `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
