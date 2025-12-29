package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

var globalAuthHandler *AuthHandler

func InitAuthHandlers(client *supabase.Client) {
	globalAuthHandler = NewAuthHandler(client)
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
	client *supabase.Client
}

func NewAuthHandler(client *supabase.Client) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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

	return c.JSON(authResponse{
		Message:   "Login realizado com sucesso",
		Token:     response.AccessToken,
		ExpiresAt: expiresAt,
		User: userResponse{
			ID:       response.User.ID,
			Username: "",
			Email:    response.User.Email,
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
			"username": req.Username,
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
			ID:       response.User.ID,
			Username: req.Username,
			Email:    response.User.Email,
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

	return c.JSON(userResponse{
		ID:       userID,
		Username: "",
		Email:    email,
	})
}
