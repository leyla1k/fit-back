package handler

import (
	"net/http"

	"github.com/folklinoff/fitness-app/internal/domain"
	middleware "github.com/folklinoff/fitness-app/internal/middleware/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

var users []domain.User
var trainers []domain.Trainer
var userID = 1
var trainerID = 1

// ResponseSuccess defines the structure for a successful response
type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseError defines the structure for an error response
type ResponseError struct {
	Error string `json:"error"`
}

// Login godoc
// @Summary Login a user or trainer
// @Description Login a user or trainer based on user_type
// @Tags auth
// @Accept json
// @Produce json
// @Param user_type path string true "User Type (user or trainer)"
// @Param credentials body domain.User true "User credentials"
// @Success 200 {object} ResponseSuccess{data=string} "token"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Router /login/{user_type} [post]
func Login(c *gin.Context) {
	userType := c.Param("user_type")

	var credentials struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data"})
		return
	}

	if userType == "user" {
		user, err := getUser(credentials.Name)
		if err != nil {
			c.JSON(http.StatusNotFound, ResponseError{Error: "no user in db"})
			return
		}
		if user.Password == credentials.Password {
			token, err := middleware.GenerateToken(uint(user.ID), userType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseError{Error: "Error generating token"})
				return
			}
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Login successful", Data: token})
			return
		}
	} else if userType == "trainer" {
		trainer, err := getTrainer(credentials.Name)
		if err != nil {
			c.JSON(http.StatusNotFound, ResponseError{Error: "no trainer in db"})
			return
		}
		if trainer.Password == credentials.Password {
			token, err := middleware.GenerateToken(uint(trainer.ID), userType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseError{Error: "Error generating token"})
				return
			}
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Login successful", Data: token})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, ResponseError{Error: "Invalid credentials"})
}

func getUser(username string) (*domain.User, error) {
	for _, user := range users {
		if user.Name == username {
			return &user, nil
		}
	}
	return nil, xerrors.Errorf("not found")
}

func getTrainer(username string) (*domain.Trainer, error) {
	for _, trainer := range trainers {
		if trainer.Name == username {
			return &trainer, nil
		}
	}
	return nil, xerrors.Errorf("not found")
}

// Register godoc
// @Summary Register a new user or trainer
// @Description Register a new user or trainer based on user_type
// @Tags auth
// @Accept json
// @Produce json
// @Param user_type path string true "User Type (user or trainer)"
// @Param user body domain.User true "User data"
// @Success 201 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Router /register/{user_type} [post]
func Register(c *gin.Context) {
	userType := c.Param("user_type")

	if userType == "user" {
		var user domain.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data"})
			return
		}
		user.ID = userID
		userID++
		users = append(users, user)
		c.JSON(http.StatusCreated, ResponseSuccess{Message: "User registered successfully"})
	} else if userType == "trainer" {
		var trainer domain.Trainer
		if err := c.ShouldBindJSON(&trainer); err != nil {
			c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data"})
			return
		}
		trainer.ID = trainerID
		trainerID++
		trainers = append(trainers, trainer)
		c.JSON(http.StatusCreated, ResponseSuccess{Message: "Trainer registered successfully"})
	}
}

// Profile godoc
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags user
// @Produce json
// @Success 200 {object} ResponseSuccess{data=domain.User}
// @Failure 401 {object} handler.ResponseError
// @Security BearerAuth
// @Router /protected/profile [get]
func Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	userType := c.MustGet("user_type").(string)

	var userProfile interface{}
	var message string

	if userType == "user" {
		for _, user := range users {
			if user.ID == userID {
				userProfile = user
				message = "User profile retrieved successfully"
				break
			}
		}
	} else if userType == "trainer" {
		for _, trainer := range trainers {
			if trainer.ID == userID {
				userProfile = trainer
				message = "Trainer profile retrieved successfully"
				break
			}
		}
	}

	if userProfile == nil {
		c.JSON(http.StatusUnauthorized, ResponseError{Error: "User not found"})
		return
	}

	response := ResponseSuccess{
		Message: message,
		Data:    userProfile,
	}

	c.JSON(http.StatusOK, response)
}
