package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"charity_system/models"

	"github.com/gin-gonic/gin"
)

// List Donations
func ListDonations(c *gin.Context) {
	// Fetch donation records, donors, and organizations
	rows, err := models.DB.Raw(`
		SELECT d.id, d.amount, d.donation_date, don.name AS donor_name, org.name AS organization_name
		FROM donations d
		JOIN donors don ON d.donor_id = don.id
		JOIN organizations org ON d.organization_id = org.id
	`).Rows()

	if err != nil {
		log.Printf("Error fetching donations: %v", err)
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to load donations.",
		})
		return
	}
	defer rows.Close()

	var donations []map[string]interface{}
	for rows.Next() {
		var id int
		var amount float64
		var donationDate, donorName, organizationName string
		if err := rows.Scan(&id, &amount, &donationDate, &donorName, &organizationName); err != nil {
			log.Printf("Error scanning donation: %v", err)
			continue
		}

		// Format the donation date (Optional: you can change the format)
		donationDateTime, err := time.Parse("2006-01-02", donationDate)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
		}
		formattedDate := donationDateTime.Format("January 2, 2006")

		donations = append(donations, map[string]interface{}{
			"id":                id,
			"amount":            amount,
			"donation_date":     formattedDate,
			"donor_name":        donorName,
			"organization_name": organizationName,
		})
	}

	// Fetch donors and organizations
	donorsRows, err := models.DB.Raw("SELECT id, name FROM donors").Rows()
	if err != nil {
		log.Printf("Error fetching donors: %v", err)
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to load donors.",
		})
		return
	}
	defer donorsRows.Close()

	var donors []map[string]interface{}
	for donorsRows.Next() {
		var id int
		var name string
		if err := donorsRows.Scan(&id, &name); err != nil {
			log.Printf("Error scanning donor: %v", err)
			continue
		}
		donors = append(donors, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	organizationsRows, err := models.DB.Raw("SELECT id, name FROM organizations").Rows()
	if err != nil {
		log.Printf("Error fetching organizations: %v", err)
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to load organizations.",
		})
		return
	}
	defer organizationsRows.Close()

	var organizations []map[string]interface{}
	for organizationsRows.Next() {
		var id int
		var name string
		if err := organizationsRows.Scan(&id, &name); err != nil {
			log.Printf("Error scanning organization: %v", err)
			continue
		}
		organizations = append(organizations, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	// Pass data to the template
	c.HTML(http.StatusOK, "donations.html", gin.H{
		"donations":     donations,
		"donors":        donors,
		"organizations": organizations,
	})
}

// Assign Donation Handler
func AssignDonation(c *gin.Context) {
	// Get the form data and validate inputs
	amountStr := c.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		log.Printf("Invalid amount: %v", err)
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Please enter a valid donation amount.",
		})
		return
	}

	donorID, err := strconv.Atoi(c.PostForm("donor_id"))
	if err != nil || donorID <= 0 {
		log.Printf("Invalid donor ID: %v", err)
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Invalid donor selected.",
		})
		return
	}

	organizationID, err := strconv.Atoi(c.PostForm("organization_id"))
	if err != nil || organizationID <= 0 {
		log.Printf("Invalid organization ID: %v", err)
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Invalid organization selected.",
		})
		return
	}

	// Insert the donation into the database
	donationDate := time.Now()
	err = models.DB.Exec("INSERT INTO donations (amount, donor_id, organization_id, donation_date) VALUES (?, ?, ?, ?)", amount, donorID, organizationID, donationDate).Error
	if err != nil {
		log.Printf("Error inserting donation: %v", err)
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to assign the donation.",
		})
		return
	}

	// Redirect to the donation list page after adding the donation
	c.Redirect(http.StatusSeeOther, "/donations")
}
