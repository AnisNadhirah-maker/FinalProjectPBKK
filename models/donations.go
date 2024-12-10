package models

import (
	"time"
)

type Donation struct {
	ID             uint      `gorm:"primaryKey"`
	Amount         float64   `gorm:"not null"`
	DonationDate   time.Time `gorm:"not null"`
	DonorID        uint      `gorm:"not null"`
	OrganizationID uint      `gorm:"not null"`
	Donor          Donor     `gorm:"foreignKey:DonorID"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}