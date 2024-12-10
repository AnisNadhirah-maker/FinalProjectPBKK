package models

type Organization struct {
    ID          uint    `gorm:"primaryKey"`
    Name        string  `gorm:"size:128;not null"`
    Description string  `gorm:"size:255"`
    Donation    float64 `gorm:"type:double"`
}

