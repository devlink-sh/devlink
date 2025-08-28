package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/devlink/devlink-api/util"
	"github.com/gin-gonic/gin"
	"github.com/openziti/edge-api/rest_management_api_client"
	"github.com/openziti/edge-api/rest_management_api_client/enrollment"
	"github.com/openziti/edge-api/rest_management_api_client/identity"
	"github.com/openziti/edge-api/rest_management_api_client/service"
	"github.com/openziti/edge-api/rest_management_api_client/service_policy"
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
	app.POST("/api/v1/create-service", createServiceHandler)

	fmt.Println("DevLink API server starting on :8080...")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createServiceHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	ctx := context.Background()

	// Check if service already exists
	listParams := &service.ListServicesParams{Context: ctx}
	listResp, err := zitiAPIClient.Service.ListServices(listParams, nil)
	if err == nil {
		for _, s := range listResp.Payload.Data {
			if s.Name != nil && *s.Name == req.Name {
				c.JSON(http.StatusOK, gin.H{"message": "service already exists", "id": s.ID})
				return
			}
		}
	} else {
		log.Printf("Error listing services: %v", err)
	}

	encRequired := false
	serviceCreate := &rest_model.ServiceCreate{
		Name:               &req.Name,
		EncryptionRequired: &encRequired,
	}
	createParams := &service.CreateServiceParams{
		Service: serviceCreate,
		Context: ctx,
	}

	createResp, err := zitiAPIClient.Service.CreateService(createParams, nil)
	if err != nil {
		switch e := err.(type) {
		case *service.CreateServiceBadRequest:
			b, jerr := json.MarshalIndent(e.Payload, "", "  ")
			if jerr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bad request creating service"})
				return
			}
			c.Data(http.StatusBadRequest, "application/json", b)
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create service", "details": err.Error()})
			return
		}
	}

	// ✅ Attach the devlink identity (from ./identity.json) to this service
	identityName := "devlink" // must match the one in provisionHandler/init.go

	// --- Bind policy ---
	bindPolicyName := req.Name + "-bind-policy"
	bindType := rest_model.DialBindBind
	semantic := rest_model.NewSemantic("AllOf")
	bindPolicy := &rest_model.ServicePolicyCreate{
		Name:          &bindPolicyName,
		Type:          &bindType,
		Semantic:      semantic,
		ServiceRoles:  []string{"#devlink-service"},
		IdentityRoles: []string{"#devlink"},
	}

	bindParams := &service_policy.CreateServicePolicyParams{
		Policy:  bindPolicy, // field is Policy in your SDK
		Context: ctx,
	}

	_, err = zitiAPIClient.ServicePolicy.CreateServicePolicy(bindParams, nil)
	if err != nil {
		if apiErr, ok := err.(*service_policy.CreateServicePolicyBadRequest); ok {
			if apiErr.Payload != nil && apiErr.Payload.Error != nil {
				log.Printf("❌ Controller rejected Bind policy: %s", apiErr.Payload.Error.Message)
			} else {
				log.Printf("❌ Controller rejected Bind policy with 400: %+v", apiErr)
			}
		} else {
			log.Printf("⚠️ Failed to attach Bind policy: %v", err)
		}
	} else {
		log.Printf("✅ Attached Bind policy: %s → %s", identityName, req.Name)
	}

	// --- Dial policy ---
	dialPolicyName := req.Name + "-dial-policy"
	dialType := rest_model.DialBindDial
	dialPolicy := &rest_model.ServicePolicyCreate{
		Name:          &dialPolicyName,
		Type:          &dialType,
		Semantic:      semantic,
		ServiceRoles:  []string{"#devlink-service"},
		IdentityRoles: []string{"#devlink"},
	}

	dialParams := &service_policy.CreateServicePolicyParams{
		Policy:  dialPolicy,
		Context: ctx,
	}

	_, err = zitiAPIClient.ServicePolicy.CreateServicePolicy(dialParams, nil)
	if err != nil {
		if apiErr, ok := err.(*service_policy.CreateServicePolicyBadRequest); ok {
			if apiErr.Payload != nil && apiErr.Payload.Error != nil {
				log.Printf("❌ Controller rejected Bind policy: %s", apiErr.Payload.Error.Message)
			} else {
				log.Printf("❌ Controller rejected Bind policy with 400: %+v", apiErr)
			}
		} else {
			log.Printf("⚠️ Failed to attach Bind policy: %v", err)
		}
	} else {
		log.Printf("✅ Attached Dial policy: %s → %s", identityName, req.Name)
	}

	c.JSON(http.StatusOK, gin.H{"message": "service created", "id": createResp.Payload.Data.ID})
}

func provisionHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if !useInvitationToken(req.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or used invitation token"})
		return
	}

	identityName := "devlink-user-" + util.GenerateHumanCode() // fixed identity for local dev
	isAdmin := false
	createIdentityParams := &identity.CreateIdentityParams{
		Identity: &rest_model.IdentityCreate{
			Enrollment: &rest_model.IdentityCreateEnrollment{
				Ott: true,
			},
			IsAdmin: &isAdmin,
			Name:    &identityName,
			Type:    rest_model.NewIdentityType("Device"),
		},
		Context: context.Background(),
	}

	resp, err := zitiAPIClient.Identity.CreateIdentity(createIdentityParams, nil)
	if err != nil {
		log.Printf("Error creating identity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identity"})
		return
	}

	identityId := resp.Payload.Data.ID

	log.Printf("Created identity with ID: %s", identityId)

	roleAttributes := rest_model.Attributes{"#devlink-sharers", "#devlink-listeners"}

	updateIdentityParams := &identity.UpdateIdentityParams{
		ID: identityId,
		Identity: &rest_model.IdentityUpdate{
			RoleAttributes: &roleAttributes,
		},
		Context: context.Background(),
	}
	_, err = zitiAPIClient.Identity.UpdateIdentity(updateIdentityParams, nil)
	if err != nil {
		log.Printf("Error assigning roles to identity %s: %v", identityId, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles to identity"})
		// return
	}
	log.Printf("Assigned roles to identity '%s'", identityName)
	

	// same enrollment logic...
	enrollmentParams := &enrollment.ListEnrollmentsParams{
		Context: context.Background(),
	}
	enrollments, err := zitiAPIClient.Enrollment.ListEnrollments(enrollmentParams, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	var enrollmentID *string
	for _, enroll := range enrollments.Payload.Data {
		if enroll.Identity != nil && enroll.Identity.ID == identityId {
			enrollmentID = enroll.ID

			break
		}
	}

	if enrollmentID == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	enrollmentDetailParams := &enrollment.DetailEnrollmentParams{
		ID:      *enrollmentID,
		Context: context.Background(),
	}

	enrollmentDetail, err := zitiAPIClient.Enrollment.DetailEnrollment(enrollmentDetailParams, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment token"})
		return
	}

	if enrollmentDetail.Payload.Data.JWT == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrollment JWT"})
		return
	}

	c.Header("Content-Type", "application/jwt")
	c.String(http.StatusOK, enrollmentDetail.Payload.Data.JWT)
	log.Printf("Provisioned identity '%s'", identityName)
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
		delete(validTokens, token)
		return true
	}
	return false
}

func initializeZitiClient() error {
	ctrlAddress := viper.GetString("ziti.controller")
	username := viper.GetString("ziti.username")
	password := viper.GetString("ziti.password")

	caPool := x509.NewCertPool()
	caCerts, err := rest_util.GetControllerWellKnownCas(ctrlAddress)
	if err != nil {
		return fmt.Errorf("failed to get controller CAs: %w", err)
	}
	for _, ca := range caCerts {
		caPool.AddCert(ca)
	}

	client, err := rest_util.NewEdgeManagementClientWithUpdb(username, password, ctrlAddress, caPool)
	if err != nil {
		return fmt.Errorf("failed to create management client: %w", err)
	}
	zitiAPIClient = client
	return nil
}
