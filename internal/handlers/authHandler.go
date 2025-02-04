package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DmitriyGiryntsev/TODO-API/internal/models"
	"github.com/DmitriyGiryntsev/TODO-API/internal/repository"
	"github.com/DmitriyGiryntsev/TODO-API/pkg/helpers"
	"github.com/DmitriyGiryntsev/TODO-API/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler defines the handler for authentication routes
type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// NewAuthHandler initializes the AuthHandler
func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user data"})
		return
	}

	userExists, err := h.UserRepo.GetUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "server error"})
		return
	}

	if userExists != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "user already exists"})
		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot hash password"})
		return
	}

	if err := h.UserRepo.CreateNewUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot create user"})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{Message: "user created successfully"})
}

// Login godoc
// @Summary Вход пользователя
// @Description Позволяет пользователю войти в систему и получить токены
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body struct{ Email string `json:"email"`; Password string `json:"password"` } true "User credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid credentials"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(creds.Email)
	if err == sql.ErrNoRows {
		log.Println("user not found")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "wrong email or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "server error"})
		return
	}

	err = utils.CheckPassword(creds.Password, user.Password)
	if err != nil {
		log.Println("wrong password")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "wrong email or password"})
		return
	}

	accessToken, refreshToken, err := helpers.GenerateAllTokens(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot generate tokens"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc
// @Summary Обновление токена
// @Description Обновляет refresh токен и выдает новый access токен
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body struct{ RefreshToken string `json:"refresh_token"` } true "Refresh token"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid refresh token"})
		return
	}

	claims, err := helpers.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "cannot create new token"})
		return
	}

	accessToken, refreshToken, err := helpers.GenerateAllTokens(claims.ID, claims.Email, claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot generate tokens"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
