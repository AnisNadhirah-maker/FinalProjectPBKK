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

	// Handle error in fetching donation records
	if err != nil {
		log.Printf("Error fetching donations: %v", err)
		// Return error message to the view
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to load donations.",
		})
		return
	}
	defer rows.Close() // Ensure to close the rows once done

	var donations []map[string]interface{}
	for rows.Next() {
		// Variables to hold the fetched data
		var id int
		var amount float64
		var donationDate, donorName, organizationName string

		// Scan each row into the respective variables
		if err := rows.Scan(&id, &amount, &donationDate, &donorName, &organizationName); err != nil {
			log.Printf("Error scanning donation: %v", err)
			continue
		}


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

	// Fetch donors from the database
	donorsRows, err := models.DB.Raw("SELECT id, name FROM donors").Rows()
	if err != nil {
		log.Printf("Error fetching donors: %v", err)
		// Return error message if fetching donors fails
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

	// Fetch organizations from the database
	organizationsRows, err := models.DB.Raw("SELECT id, name FROM organizations").Rows()
	if err != nil {
		log.Printf("Error fetching organizations: %v", err)
		// Return error message if fetching organizations fails
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

	// Pass data to the HTML template for rendering the page
	c.HTML(http.StatusOK, "donations.html", gin.H{
		"donations":     donations,     
		"donors":        donors,        
		"organizations": organizations, 
	})
}

// Assign Donation Handler
func AssignDonation(c *gin.Context) {
	// Get the donation amount from the form and validate it
	amountStr := c.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		log.Printf("Invalid amount: %v", err)
		// Return an error if the amount is invalid
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Please enter a valid donation amount.",
		})
		return
	}

	// Get the donor ID from the form and validate it
	donorID, err := strconv.Atoi(c.PostForm("donor_id"))
	if err != nil || donorID <= 0 {
		log.Printf("Invalid donor ID: %v", err)
		// Return an error if the donor ID is invalid
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Invalid donor selected.",
		})
		return
	}

	// Get the organization ID from the form and validate it
	organizationID, err := strconv.Atoi(c.PostForm("organization_id"))
	if err != nil || organizationID <= 0 {
		log.Printf("Invalid organization ID: %v", err)
		// Return an error if the organization ID is invalid
		c.HTML(http.StatusBadRequest, "donations.html", gin.H{
			"error": "Invalid organization selected.",
		})
		return
	}

	// Insert the donation record into the database
	donationDate := time.Now() 
	err = models.DB.Exec("INSERT INTO donations (amount, donor_id, organization_id, donation_date) VALUES (?, ?, ?, ?)", amount, donorID, organizationID, donationDate).Error
	if err != nil {
		log.Printf("Error inserting donation: %v", err)
		c.HTML(http.StatusInternalServerError, "donations.html", gin.H{
			"error": "Failed to assign the donation.",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/donations")
}
