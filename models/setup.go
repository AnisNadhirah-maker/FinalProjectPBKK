package models

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:ca23010@tcp(127.0.0.1:3306)/charity_system?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	DB = database
}

// DBMigrate runs the database migrations
func DBMigrate() {
	DB.AutoMigrate(&Donor{}, &Organization{},&Donation{},
	)	
}

