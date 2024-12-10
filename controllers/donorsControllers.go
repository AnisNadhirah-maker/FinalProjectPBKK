package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"charity_system/models"
	"github.com/gin-gonic/gin"
)

// ListDonors retrieves all donors and renders them in HTML
func ListDonors(c *gin.Context) {
	var donors []models.Donor
	if err := models.DB.Find(&donors).Error; err != nil {
		log.Println("Error retrieving donors: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donors."})
		return
	}

	donorData := make([]map[string]interface{}, len(donors))
	for i, donor := range donors {
		donorData[i] = map[string]interface{}{
			"id":    donor.ID,
			"name":  donor.Name,
			"email": donor.Email,
		}
	}

	c.HTML(http.StatusOK, "donors.html", gin.H{"donors": donorData})
}

// AddDonor handles adding a new donor to the database
func AddDonor(c *gin.Context) {
	name := c.PostForm("name")
	Email := c.PostForm("email")

	if name == "" || Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required fields."})
		return
	}

	donor := models.Donor{
		Name:      name,
		Email:     Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.DB.Create(&donor).Error; err != nil {
		log.Println("Error inserting donor into database: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add donor."})
		return
	}

	c.Redirect(http.StatusSeeOther, "/donors")
}

// DeleteDonor handles deleting donor records
func DeleteDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID."})
		return
	}

	if err := models.DB.Delete(&models.Donor{}, id).Error; err != nil {
		log.Println("Error deleting donor: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete donor."})
		return
	}

	c.Redirect(http.StatusSeeOther, "/donors")
}
