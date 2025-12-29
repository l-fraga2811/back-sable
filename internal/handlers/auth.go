package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/repository"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

var globalAuthHandler *AuthHandler

func InitAuthHandlers(handler *AuthHandler) {
	globalAuthHandler = handler
}

func SignIn(c fiber.Ctx) error {
	if globalAuthHandler == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Auth handler not initialized",
		})
	}
	return globalAuthHandler.Login(c)
}

func SignUp(c fiber.Ctx) error {
	if globalAuthHandler == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Auth handler not initialized",
		})
	}
	return globalAuthHandler.Register(c)
}

func GetProfile(c fiber.Ctx) error {
	if globalAuthHandler == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Auth handler not initialized",
		})
	}
	return globalAuthHandler.GetProfile(c)
}

type AuthHandler struct {
	client      *supabase.Client
	profileRepo repository.ProfileRepository
}

func NewAuthHandler(client *supabase.Client) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func NewAuthHandlerWithProfileRepo(client *supabase.Client, profileRepo repository.ProfileRepository) *AuthHandler {
	return &AuthHandler{
		client:      client,
		profileRepo: profileRepo,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	ProfileUrl string `json:"profileUrl"`
}

type userResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	ProfileUrl string `json:"profileUrl"`
}

type authResponse struct {
	Message   string       `json:"message"`
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expiresAt"`
	User      userResponse `json:"user"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.client.SignIn(c.Context(), supabase.SignInCredentials{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication failed: " + err.Error(),
		})
	}

	expiresAt := ""
	if response.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
	}

	username := ""
	if response.User.UserMetadata != nil {
		fmt.Println(response)
		if usernameValue, ok := response.User.UserMetadata["username"]; ok {
			if usernameStr, ok := usernameValue.(string); ok {
				username = usernameStr
			}
		}
	}

	return c.JSON(authResponse{
		Message:   "Login realizado com sucesso",
		Token:     response.AccessToken,
		ExpiresAt: expiresAt,
		User: userResponse{
			ID:         response.User.ID,
			Username:   username,
			Email:      response.User.Email,
			Phone:      response.User.UserMetadata["phone"].(string),
			ProfileUrl: response.User.UserMetadata["profile_url"].(string),
		},
	})
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req registerRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.client.SignUp(c.Context(), supabase.SignUpCredentials{
		Email:    req.Email,
		Password: req.Password,
		Data: map[string]interface{}{
			"username":    req.Username,
			"phone":       req.Phone,
			"profile_url": req.ProfileUrl,
		},
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Registration failed: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Usuário criado com sucesso",
		"user": userResponse{
			ID:         response.User.ID,
			Username:   req.Username,
			Email:      response.User.Email,
			Phone:      req.Phone,
			ProfileUrl: req.ProfileUrl,
		},
	})
}

func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	email, _ := c.Locals("email").(string)
	username, _ := c.Locals("username").(string)

	if h.profileRepo == nil {
		return c.Status(fiber.StatusOK).JSON(userResponse{
			ID:         userID,
			Username:   username,
			Email:      email,
			Phone:      "",
			ProfileUrl: "",
		})
	}

	profile, err := h.profileRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load profile",
		})
	}

	if profile == nil {
		return c.Status(fiber.StatusOK).JSON(userResponse{
			ID:         userID,
			Username:   username,
			Email:      email,
			Phone:      "",
			ProfileUrl: "",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userResponse{
		ID:         userID,
		Username:   profile.Username,
		Email:      email,
		Phone:      profile.Phone,
		ProfileUrl: profile.ProfileUrl,
	})
}
