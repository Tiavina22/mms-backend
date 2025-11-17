package services

import (
	"errors"
	"strings"

	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// SignupRequest represents signup request data
type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
	Language string `json:"language"`
}

// LoginRequest represents login request data
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Can be email, username, or phone
	Password   string `json:"password" binding:"required"`
}

// AuthResponse represents auth response with token and user data
type AuthResponse struct {
	Token string            `json:"token"`
	User  models.PublicUser `json:"user"`
}

// Signup registers a new user
func (s *AuthService) Signup(req SignupRequest) (*AuthResponse, error) {
	// Validate input
	if err := utils.ValidateUsername(req.Username); err != nil {
		return nil, err
	}
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			return nil, err
		}
	}

	// Check if user already exists
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already taken")
	}
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}
	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(req.Phone); err == nil {
			return nil, errors.New("phone already registered")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Determine language
	language := req.Language
	if language == "" && req.Phone != "" {
		language = utils.GetLanguageFromPhone(req.Phone)
	}
	if language == "" {
		language = "en"
	}

	// Create user
	user := &models.User{
		Username: utils.SanitizeString(req.Username),
		Email:    strings.ToLower(utils.SanitizeString(req.Email)),
		Phone:    utils.SanitizeString(req.Phone),
		Password: hashedPassword,
		Language: language,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user.ToPublicUser(),
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Find user by identifier (email, username, or phone)
	var user *models.User
	var err error

	identifier := strings.ToLower(utils.SanitizeString(req.Identifier))

	// Try to find by email first
	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.FindByEmail(identifier)
	} else if strings.HasPrefix(identifier, "+") {
		// Try phone
		user, err = s.userRepo.FindByPhone(identifier)
	} else {
		// Try username
		user, err = s.userRepo.FindByUsername(identifier)
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Update online status
	_ = s.userRepo.UpdateOnlineStatus(user.ID, true)

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user.ToPublicUser(),
	}, nil
}

// Logout logs out a user
func (s *AuthService) Logout(userID string) error {
	// Update online status
	// Parse userID string to UUID (assuming it's already validated in middleware)
	return nil
}

// ValidateToken validates a JWT token and returns user info
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

