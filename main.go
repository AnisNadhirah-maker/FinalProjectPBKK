package main

import (
	"log"
	"net/http"

	"charity_system/controllers"
	"charity_system/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Static and HTML files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Connect to the database and apply migrations
	models.ConnectDatabase()
	models.DB = models.DB.Debug()
	models.DBMigrate()
	log.Println("Database migration completed successfully")

	// Routes
	r.GET("/", homePage)

	// Donor routes
	r.GET("/donors", controllers.ListDonors) 
	r.POST("/donors/add", controllers.AddDonor)
	r.DELETE("/donors/:id", controllers.DeleteDonor)  // Changed to DELETE

	// Organization routes
	r.GET("/organizations", controllers.ListOrganizations)           
	r.POST("/organizations/add", controllers.AddOrganization)        
	r.DELETE("/organizations/:id", controllers.DeleteOrganization)  // Changed to DELETE

	// Donation routes
	r.GET("/donations", controllers.ListDonations)                    
	r.POST("/donations/assign", controllers.AssignDonation)

	// Start the server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
