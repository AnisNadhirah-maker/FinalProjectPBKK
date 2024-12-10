package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"charity_system/models"
)

// ListDonations retrieves and displays all donations along with associated donors and organizations
func ListDonations(c *gin.Context) {
	var donations []models.Donation

	// Fetch donations with associated donors and organizations
	if err := models.DB.Preload("Donor").Preload("Organization").Find(&donations).Error; err != nil {
		log.Println("Error retrieving donations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations."})
		return
	}

	// Prepare donations for rendering
	donationData := make([]map[string]interface{}, len(donations))
	for i, donation := range donations {
		donationData[i] = map[string]interface{}{
			"id":            donation.ID,
			"amount":        donation.Amount,
			"donation_date": donation.DonationDate.Format("2006-01-02"),
			"donor_name":    donation.Donor.Name,
			"organization":  donation.Organization.Name,
		}
	}

	// Render the donations.html template with the data
	c.HTML(http.StatusOK, "donations.html", gin.H{"donations": donationData})
}

func AssignDonation(c *gin.Context) {
    // Parse donation_id, donor_id, and organization_id from form data
    donationID, err := strconv.Atoi(c.PostForm("donation_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
        return
    }
    donorID, err := strconv.Atoi(c.PostForm("donor_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
        return
    }
    organizationID, err := strconv.Atoi(c.PostForm("organization_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
        return
    }

    // Fetch the donation record
    var donation models.Donation
    if err := models.DB.First(&donation, donationID).Error; err != nil {
        log.Println("Error retrieving donation record:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Donation not found"})
        return
    }

    // Update the donation record
    donation.DonorID = uint(donorID)
    donation.OrganizationID = uint(organizationID)
    if err := models.DB.Save(&donation).Error; err != nil {
        log.Println("Error updating donation:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign donation"})
        return
    }

    // Update donor status
    if err := models.DB.Model(&models.Donor{}).Where("id = ?", donorID).Update("status", "Assigned").Error; err != nil {
        log.Println("Error updating donor status:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update donor status"})
        return
    }

    // Redirect to donations page
    c.Redirect(http.StatusSeeOther, "/donations")
}
