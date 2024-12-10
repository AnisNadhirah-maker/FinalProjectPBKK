package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"charity_system/models"
)

// ListOrganizations retrieves and displays all organizations
func ListOrganizations(c *gin.Context) {
	var organizations []models.Organization
	if err := models.DB.Find(&organizations).Error; err != nil {
		log.Println("Error retrieving organizations: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations."})
		return
	}

	organizationData := make([]map[string]interface{}, len(organizations))
	for i, organization := range organizations {
		organizationData[i] = map[string]interface{}{
			"id":          organization.ID,
			"name":        organization.Name,
			"description": organization.Description,
			"donation":    organization.Donation,
		}
	}
	c.HTML(http.StatusOK, "organizations.html", gin.H{"organizations": organizationData})
}

// AddOrganization handles adding a new organization to the database
func AddOrganization(c *gin.Context) {
	// Retrieve form data for organization
	name := c.PostForm("name")
	description := c.PostForm("description")
	donation := c.PostForm("donation")

	// Validate required fields
	if name == "" || description == "" || donation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, description, and donation are required fields."})
		return
	}

	// Convert donation from string to float64
	donationValue, err := strconv.ParseFloat(donation, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation value."})
		return
	}

	// Create a new organization record
	org := models.Organization{
		Name:        name,
		Description: description,
		Donation:    donationValue,
	}

	// Insert the new organization into the database
	if err := models.DB.Create(&org).Error; err != nil {
		log.Println("Error inserting organization into database: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add organization."})
		return
	}

	// Redirect to the organizations page
	c.Redirect(http.StatusSeeOther, "/organizations")
}


// DeleteOrganization handles deleting an organization record
func DeleteOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID."})
		return
	}

	if err := models.DB.Delete(&models.Organization{}, id).Error; err != nil {
		log.Println("Error deleting organization: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization."})
		return
	}
	c.Redirect(http.StatusSeeOther, "/organizations")
}
