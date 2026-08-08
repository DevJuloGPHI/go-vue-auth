package auth

import (
	"errors"
	"net/http"

	"github.com/DevJuloGPHI/go-vue-auth/backend/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService user.Service
}

func NewHandler(userService user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid registration data",
		})
		return
	}

	createdUser, err := h.userService.Register(
		request.Name,
		request.Email,
		request.Password,
	)

	if errors.Is(err, user.ErrEmailAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"id":    createdUser.ID,
			"name":  createdUser.Name,
			"email": createdUser.Email,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login data",
		})
		return
	}

	loggedInUser, err := h.userService.Login(
		request.Email,
		request.Password,
	)

	if errors.Is(err, user.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    loggedInUser.ID,
			"name":  loggedInUser.Name,
			"email": loggedInUser.Email,
		},
	})
}
