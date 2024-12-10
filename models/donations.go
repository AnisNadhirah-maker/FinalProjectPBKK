package models

import (
	"time"
)

type Donation struct {
    ID             uint        `gorm:"primaryKey"`
    Amount         float64     `gorm:"not null"`
    DonationDate   time.Time   `gorm:"not null"`
    OrganizationID uint        `gorm:"not null"`
    DonorID        uint        `gorm:"not null"`
    Organization   Organization `gorm:"foreignKey:OrganizationID"`
    Donor          Donor        `gorm:"foreignKey:DonorID"`
}