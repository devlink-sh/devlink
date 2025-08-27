package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openziti/edge-api/rest_management_api_client"
	"github.com/openziti/edge-api/rest_management_api_client/enrollment"
	"github.com/openziti/edge-api/rest_management_api_client/identity"
	"github.com/openziti/edge-api/rest_model"
	"github.com/openziti/edge-api/rest_util"
	"github.com/spf13/viper"
)

var (
	validTokens   = make(map[string]bool)
	tokensMutex   = &sync.Mutex{}
	zitiAPIClient *rest_management_api_client.ZitiEdgeManagement
)

func main() {

	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := loadInvitationTokens("/home/rishi/Desktop/devlink/invites.txt"); err != nil {
		log.Fatalf("Failed to load invitation tokens: %v", err)
	}

	if err := initializeZitiClient(); err != nil {
		log.Fatalf("Failed to initialize Ziti client: %v", err)
	}

	app := gin.New()

	app.POST("/api/v1/provision", provisionHandler)

	fmt.Println("DevLink API server starting on :8080...")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func provisionHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	// Decode the incoming request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the invitation token is valid and use it up.
	if !useInvitationToken(req.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or used invitation token"})
		return
	}

	// The token is valid, so create a new identity in OpenZiti.
	// We'll give the identity a unique name.
	identityName := "devlink-user-" + uuid.New().String()

	// This is the API call to create the identity.
	// We specify `enrollment: { "ott": true }` to get a one-time token.
	isAdmin := false
	createIdentityParams := &identity.CreateIdentityParams{
		Identity: &rest_model.IdentityCreate{
			Enrollment: &rest_model.IdentityCreateEnrollment{
				Ott: true,
			},
			IsAdmin: &isAdmin,
			Name:    &identityName,
			Type:    rest_model.NewIdentityType("Device"), // Or "User"
		},
		Context: context.Background(),
	}

	resp, err := zitiAPIClient.Identity.CreateIdentity(createIdentityParams, nil)
	if err != nil {
		log.Printf("Error creating identity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identity"})
		return
	}

	log.Printf("Created identity with ID: %s", resp.Payload.Data.ID)

	// List all enrollments and find the one for our identity
	enrollmentParams := &enrollment.ListEnrollmentsParams{
		Context: context.Background(),
	}
	enrollments, err := zitiAPIClient.Enrollment.ListEnrollments(enrollmentParams, nil)
	if err != nil {
		log.Printf("Error retrieving enrollments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	// Find the enrollment for our identity
	var enrollmentID *string
	for _, enroll := range enrollments.Payload.Data {
		log.Printf("Checking enrollment ID: %s, Identity: %v, Token: %v",
			*enroll.ID,
			enroll.Identity,
			enroll.Token != nil)

		if enroll.Identity != nil && enroll.Identity.ID == resp.Payload.Data.ID {
			enrollmentID = enroll.ID
			log.Printf("Found matching enrollment ID for identity %s: %s", resp.Payload.Data.ID, *enrollmentID)
			break
		}
	}

	if enrollmentID == nil {
		log.Printf("No enrollment found for identity %s", resp.Payload.Data.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	// Now get the enrollment details to fetch the actual JWT
	enrollmentDetailParams := &enrollment.DetailEnrollmentParams{
		ID:      *enrollmentID,
		Context: context.Background(),
	}

	enrollmentDetail, err := zitiAPIClient.Enrollment.DetailEnrollment(enrollmentDetailParams, nil)
	if err != nil {
		log.Printf("Error retrieving enrollment detail for ID %s: %v", *enrollmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	if enrollmentDetail.Payload.Data.JWT == "" {
		log.Printf("No JWT found in enrollment detail for ID %s", *enrollmentID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment JWT"})
		return
	}

	log.Printf("Sending enrollment JWT: %s", enrollmentDetail.Payload.Data.JWT)
	c.Header("Content-Type", "application/jwt")
	c.String(http.StatusOK, enrollmentDetail.Payload.Data.JWT)

	log.Printf("Successfully provisioned identity '%s'", identityName)
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func loadInvitationTokens(filename string) error {
	tokensMutex.Lock()
	defer tokensMutex.Unlock()

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		token := strings.TrimSpace(line)
		if token != "" {
			validTokens[token] = true
		}
	}
	fmt.Printf("Loaded %d invitation tokens\n", len(validTokens))
	return nil
}

func useInvitationToken(token string) bool {
	tokensMutex.Lock()
	defer tokensMutex.Unlock()

	if validTokens[token] {
		delete(validTokens, token) // Token is now used up
		return true
	}
	return false
}

func initializeZitiClient() error {
	ctrlAddress := viper.GetString("ziti.controller")
	username := viper.GetString("ziti.username")
	password := viper.GetString("ziti.password")

	// Trust the controller's CA
	caPool := x509.NewCertPool()
	caCerts, err := rest_util.GetControllerWellKnownCas(ctrlAddress)
	if err != nil {
		return fmt.Errorf("failed to get controller CAs: %w", err)
	}
	for _, ca := range caCerts {
		caPool.AddCert(ca)
	}

	// Authenticate and create a management client
	client, err := rest_util.NewEdgeManagementClientWithUpdb(username, password, ctrlAddress, caPool)
	if err != nil {
		return fmt.Errorf("failed to create management client: %w", err)
	}
	zitiAPIClient = client
	return nil
}
