package models

import (

	"time"
)

// Donor struct represents the donors table
type Donor struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:128"`
    Email     string    `gorm:"size:128"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

