package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mms-backend/config"
	"mms-backend/controllers"
	"mms-backend/middleware"
	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/routes"
	"mms-backend/services"
	"mms-backend/websocket"
)

var (
	router     *gin.Engine
	db         *gorm.DB
	testServer *httptest.Server
)

// Test user credentials
var (
	testUserAlice = map[string]interface{}{
		"username": "alice_test",
		"email":    "alice_test@example.com",
		"password": "Alice1234!",
		"phone":    "+261340000001",
		"language": "fr",
	}
	testUserBob = map[string]interface{}{
		"username": "bob_test",
		"email":    "bob_test@example.com",
		"password": "Bob1234!",
		"phone":    "+1234567890",
		"language": "en",
	}

	aliceToken string
	aliceID    string
	bobToken   string
	bobID      string
)

func TestMain(m *testing.M) {
	// Setup
	setupTestEnvironment()
	setupTestDatabase()
	setupTestRouter()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestDatabase()

	os.Exit(code)
}

func setupTestEnvironment() {
	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "4dmin")
	os.Setenv("DB_NAME", "mms_test")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("PORT", "8081")
	os.Setenv("ENV", "test")
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-for-integration-tests")
	os.Setenv("JWT_EXPIRY", "24h")
	os.Setenv("ENCRYPTION_KEY", "test-encryption-key-32-bytes!!")

	// Load config
	config.LoadConfig()
}

func setupTestDatabase() {
	cfg := config.AppConfig

	// Connect to PostgreSQL to create test database
	dsn := "host=" + cfg.Database.Host + " port=" + cfg.Database.Port +
		" user=" + cfg.Database.User + " password=" + cfg.Database.Password +
		" dbname=postgres sslmode=" + cfg.Database.SSLMode

	adminDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}

	// Drop and recreate test database
	sqlDB, _ := adminDB.DB()
	sqlDB.Exec("DROP DATABASE IF EXISTS mms_test")
	sqlDB.Exec("CREATE DATABASE mms_test")
	sqlDB.Close()

	// Connect to test database
	testDSN := cfg.Database.GetDSN()
	db, err = gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Run migrations
	db.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Group{},
		&models.GroupMember{},
		&models.GroupMessage{},
		&models.Notification{},
	)
}

func setupTestRouter() {
	gin.SetMode(gin.TestMode)
	router = gin.New()
	router.Use(middleware.CORSMiddleware())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	groupMessageRepo := repositories.NewGroupMessageRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)

	// Initialize services
	pushService := services.NewPushService(config.AppConfig)
	authService := services.NewAuthService(userRepo)
	messageService := services.NewMessageService(messageRepo, userRepo, notificationRepo, pushService)
	groupService := services.NewGroupService(groupRepo, groupMessageRepo, userRepo, notificationRepo, pushService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userRepo)
	messageController := controllers.NewMessageController(messageService)
	groupController := controllers.NewGroupController(groupService)

	// Initialize WebSocket
	hub := websocket.NewHub()
	wsHandler := websocket.NewHandler(hub)

	// Setup routes
	routes.SetupRoutes(router, authController, userController, messageController, groupController, wsHandler)
}

func cleanupTestDatabase() {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

// Helper functions
func makeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func parseResponse(w *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), target)
}

// ========================================
// AUTHENTICATION TESTS
// ========================================

func TestHealthCheck(t *testing.T) {
	w := makeRequest("GET", "/health", nil, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	assert.Equal(t, "ok", response["status"])
	assert.NotEmpty(t, response["message"])
	
	t.Log("✓ Health check passed")
}

func TestSignupAlice(t *testing.T) {
	w := makeRequest("POST", "/api/v1/auth/signup", testUserAlice, "")
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	assert.NotNil(t, response["data"])
	data := response["data"].(map[string]interface{})
	
	aliceToken = data["token"].(string)
	user := data["user"].(map[string]interface{})
	aliceID = user["id"].(string)
	
	assert.NotEmpty(t, aliceToken)
	assert.NotEmpty(t, aliceID)
	assert.Equal(t, "alice_test", user["username"])
	assert.Equal(t, "fr", user["language"])
	
	t.Logf("✓ Alice signed up successfully - ID: %s", aliceID)
}

func TestSignupBob(t *testing.T) {
	w := makeRequest("POST", "/api/v1/auth/signup", testUserBob, "")
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	bobToken = data["token"].(string)
	user := data["user"].(map[string]interface{})
	bobID = user["id"].(string)
	
	assert.NotEmpty(t, bobToken)
	assert.NotEmpty(t, bobID)
	assert.Equal(t, "bob_test", user["username"])
	assert.Equal(t, "en", user["language"])
	
	t.Logf("✓ Bob signed up successfully - ID: %s", bobID)
}

func TestSignupDuplicateEmail(t *testing.T) {
	w := makeRequest("POST", "/api/v1/auth/signup", testUserAlice, "")
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	assert.Contains(t, response["error"], "already")
	
	t.Log("✓ Duplicate email correctly rejected")
}

func TestLoginWithEmail(t *testing.T) {
	loginData := map[string]string{
		"identifier": testUserAlice["email"].(string),
		"password":   testUserAlice["password"].(string),
	}
	
	w := makeRequest("POST", "/api/v1/auth/login", loginData, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	token := data["token"].(string)
	
	assert.NotEmpty(t, token)
	
	t.Log("✓ Login with email successful")
}

func TestLoginWithUsername(t *testing.T) {
	loginData := map[string]string{
		"identifier": testUserBob["username"].(string),
		"password":   testUserBob["password"].(string),
	}
	
	w := makeRequest("POST", "/api/v1/auth/login", loginData, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	t.Log("✓ Login with username successful")
}

func TestLoginInvalidCredentials(t *testing.T) {
	loginData := map[string]string{
		"identifier": testUserAlice["email"].(string),
		"password":   "WrongPassword123!",
	}
	
	w := makeRequest("POST", "/api/v1/auth/login", loginData, "")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	t.Log("✓ Invalid credentials correctly rejected")
}

func TestGetMe(t *testing.T) {
	w := makeRequest("GET", "/api/v1/auth/me", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "alice_test", data["username"])
	assert.Equal(t, aliceID, data["id"])
	
	t.Log("✓ Get current user successful")
}

func TestGetMeWithoutToken(t *testing.T) {
	w := makeRequest("GET", "/api/v1/auth/me", nil, "")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	t.Log("✓ Unauthorized access correctly rejected")
}

// ========================================
// USER TESTS
// ========================================

func TestListUsers(t *testing.T) {
	w := makeRequest("GET", "/api/v1/users?limit=10", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2) // At least Alice and Bob
	
	t.Logf("✓ Listed %d users", len(data))
}

func TestSearchUsers(t *testing.T) {
	w := makeRequest("GET", "/api/v1/users/search?q=bob&limit=5", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
	
	// Verify Bob is in results
	found := false
	for _, user := range data {
		u := user.(map[string]interface{})
		if u["username"] == "bob_test" {
			found = true
			break
		}
	}
	assert.True(t, found, "Bob should be found in search results")
	
	t.Log("✓ User search successful")
}

func TestGetUserByID(t *testing.T) {
	w := makeRequest("GET", "/api/v1/users/"+bobID, nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "bob_test", data["username"])
	
	t.Log("✓ Get user by ID successful")
}

// ========================================
// MESSAGE TESTS
// ========================================

func TestSendMessage(t *testing.T) {
	messageData := map[string]interface{}{
		"receiver_id": bobID,
		"content":     "Hello Bob! This is a test message.",
	}
	
	w := makeRequest("POST", "/api/v1/messages", messageData, aliceToken)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, aliceID, data["sender_id"])
	assert.Equal(t, bobID, data["receiver_id"])
	assert.Equal(t, "Hello Bob! This is a test message.", data["content"])
	
	t.Log("✓ Message sent successfully")
}

func TestSendMultipleMessages(t *testing.T) {
	messages := []string{
		"First message",
		"Second message",
		"Third message",
	}
	
	for _, msg := range messages {
		messageData := map[string]interface{}{
			"receiver_id": bobID,
			"content":     msg,
		}
		
		w := makeRequest("POST", "/api/v1/messages", messageData, aliceToken)
		assert.Equal(t, http.StatusCreated, w.Code)
	}
	
	t.Log("✓ Multiple messages sent successfully")
}

func TestGetConversation(t *testing.T) {
	w := makeRequest("GET", "/api/v1/messages/conversation/"+bobID+"?limit=50", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 3) // At least 3 messages
	
	t.Logf("✓ Retrieved %d messages from conversation", len(data))
}

func TestGetUnreadCount(t *testing.T) {
	w := makeRequest("GET", "/api/v1/messages/unread/count", nil, bobToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	count := response["count"].(float64)
	assert.GreaterOrEqual(t, count, float64(4)) // At least 4 unread messages from Alice
	
	t.Logf("✓ Bob has %d unread messages", int(count))
}

func TestMarkAsRead(t *testing.T) {
	w := makeRequest("PUT", "/api/v1/messages/read/"+aliceID, nil, bobToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify count is now 0
	w2 := makeRequest("GET", "/api/v1/messages/unread/count", nil, bobToken)
	var response map[string]interface{}
	parseResponse(w2, &response)
	
	count := response["count"].(float64)
	assert.Equal(t, float64(0), count)
	
	t.Log("✓ Messages marked as read successfully")
}

func TestGetRecentConversations(t *testing.T) {
	w := makeRequest("GET", "/api/v1/messages/conversations?limit=10", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
	
	t.Logf("✓ Retrieved %d recent conversations", len(data))
}

// ========================================
// GROUP TESTS
// ========================================

var testGroupID string

func TestCreateGroup(t *testing.T) {
	groupData := map[string]interface{}{
		"name":        "Test Group",
		"description": "A test group for integration testing",
		"type":        "private",
		"member_ids":  []string{bobID},
	}
	
	w := makeRequest("POST", "/api/v1/groups", groupData, aliceToken)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	testGroupID = data["id"].(string)
	
	assert.NotEmpty(t, testGroupID)
	assert.Equal(t, "Test Group", data["name"])
	assert.Equal(t, "private", data["type"])
	
	t.Logf("✓ Group created successfully - ID: %s", testGroupID)
}

func TestGetUserGroups(t *testing.T) {
	w := makeRequest("GET", "/api/v1/groups/my", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
	
	t.Logf("✓ Alice has %d groups", len(data))
}

func TestGetGroup(t *testing.T) {
	w := makeRequest("GET", "/api/v1/groups/"+testGroupID, nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Group", data["name"])
	
	t.Log("✓ Get group successful")
}

func TestSendGroupMessage(t *testing.T) {
	messageData := map[string]interface{}{
		"group_id": testGroupID,
		"content":  "Hello everyone in the group!",
	}
	
	w := makeRequest("POST", "/api/v1/groups/messages", messageData, aliceToken)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, testGroupID, data["group_id"])
	assert.Equal(t, "Hello everyone in the group!", data["content"])
	
	t.Log("✓ Group message sent successfully")
}

func TestBobReplyInGroup(t *testing.T) {
	messageData := map[string]interface{}{
		"group_id": testGroupID,
		"content":  "Thanks Alice! Great to be here!",
	}
	
	w := makeRequest("POST", "/api/v1/groups/messages", messageData, bobToken)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	t.Log("✓ Bob replied in group successfully")
}

func TestGetGroupMessages(t *testing.T) {
	w := makeRequest("GET", "/api/v1/groups/"+testGroupID+"/messages?limit=50", nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2) // At least 2 messages
	
	t.Logf("✓ Retrieved %d group messages", len(data))
}

func TestDeleteGroup(t *testing.T) {
	w := makeRequest("DELETE", "/api/v1/groups/"+testGroupID, nil, aliceToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify group is deleted - members and messages should be gone
	var memberCount int64
	db.Model(&models.GroupMember{}).Where("group_id = ?", testGroupID).Count(&memberCount)
	assert.Equal(t, int64(0), memberCount, "Group members should be deleted")
	
	var messageCount int64
	db.Model(&models.GroupMessage{}).Where("group_id = ?", testGroupID).Count(&messageCount)
	assert.Equal(t, int64(0), messageCount, "Group messages should be deleted")
	
	t.Log("✓ Group deleted successfully with all members and messages")
}

// ========================================
// SECURITY TESTS
// ========================================

func TestUnauthorizedAccess(t *testing.T) {
	// Test GET endpoints
	getEndpoints := []string{
		"/api/v1/auth/me",
		"/api/v1/users",
		"/api/v1/messages/conversations",
		"/api/v1/groups/my",
	}
	
	for _, endpoint := range getEndpoints {
		w := makeRequest("GET", endpoint, nil, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code, "GET "+endpoint+" should reject unauthorized access")
	}
	
	// Test POST endpoints
	postEndpoints := []struct {
		path string
		body interface{}
	}{
		{"/api/v1/messages", map[string]string{"receiver_id": "test", "content": "test"}},
		{"/api/v1/groups", map[string]string{"name": "test"}},
	}
	
	for _, endpoint := range postEndpoints {
		w := makeRequest("POST", endpoint.path, endpoint.body, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code, "POST "+endpoint.path+" should reject unauthorized access")
	}
	
	t.Log("✓ All protected endpoints correctly reject unauthorized access")
}

func TestInvalidToken(t *testing.T) {
	w := makeRequest("GET", "/api/v1/auth/me", nil, "invalid.token.here")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	t.Log("✓ Invalid token correctly rejected")
}

// ========================================
// ENCRYPTION TESTS
// ========================================

func TestMessageEncryption(t *testing.T) {
	// Send a message
	messageData := map[string]interface{}{
		"receiver_id": bobID,
		"content":     "This message should be encrypted in the database",
	}
	
	w := makeRequest("POST", "/api/v1/messages", messageData, aliceToken)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	parseResponse(w, &response)
	data := response["data"].(map[string]interface{})
	messageID := data["id"].(string)
	
	// Check in database
	var msg models.Message
	db.First(&msg, "id = ?", messageID)
	
	// Content in DB should be encrypted (base64, different from original)
	assert.NotEqual(t, "This message should be encrypted in the database", msg.Content)
	assert.NotEmpty(t, msg.Content)
	
	t.Log("✓ Message is encrypted in database")
}

