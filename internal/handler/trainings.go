package handler

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/folklinoff/fitness-app/internal/domain"
	"github.com/gin-gonic/gin"
)

var trainings []domain.Training
var trainingID = 1

// CreateTraining godoc
// @Summary Create a new training session
// @Description Create a new training session (only for trainers)
// @Tags training
// @Accept json
// @Produce json
// @Param training body domain.Training true "Training data"
// @Success 201 {object} ResponseSuccess{data=domain.Training}
// @Failure 400 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training [post]
func CreateTraining(c *gin.Context) {
	trainerID := c.MustGet("user_id").(float64)

	var training domain.Training
	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data: " + err.Error()})
		return
	}

	training.ID = trainingID
	training.TrainerID = int(trainerID)
	trainingID++
	trainings = append(trainings, training)

	for i, trainer := range trainers {
		if trainer.ID == int(trainerID) {
			trainers[i].Trainings = append(trainers[i].Trainings, training.ID)
			break
		}
	}

	c.JSON(http.StatusCreated, ResponseSuccess{Message: "Training created successfully", Data: training})
}

// RegisterUserForTraining godoc
// @Summary Register a user for a training session
// @Description Register a user for a specific training session
// @Tags training
// @Accept json
// @Produce json
// @Param training_id path int true "Training ID"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training/{training_id}/register [post]
func RegisterUserForTraining(c *gin.Context) {
	userID := c.MustGet("user_id").(float64)
	trainingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid training ID: " + err.Error()})
		return
	}

	for i, training := range trainings {
		if training.ID == trainingID {
			for _, user := range training.Users {
				if user == int(userID) {
					c.JSON(http.StatusConflict, ResponseError{Error: "User already registered for this training"})
					return
				}
			}
			trainings[i].Users = append(trainings[i].Users, int(userID))

			for j, user := range users {
				if user.ID == int(userID) {
					users[j].Trainings = append(users[j].Trainings, trainingID)
					break
				}
			}

			c.JSON(http.StatusOK, ResponseSuccess{Message: "User registered for training"})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "Training not found with ID " + strconv.Itoa(trainingID)})
}

// GetTrainingByID godoc
// @Summary Get a training session by ID
// @Description Get a training session by ID
// @Tags training
// @Produce json
// @Param id path int true "Training ID"
// @Success 200 {object} ResponseSuccess{data=domain.Training}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training/{id} [get]
func GetTrainingByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid training ID: " + err.Error()})
		return
	}

	for _, training := range trainings {
		if training.ID == id {
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Training found", Data: training})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "Training not found with ID " + strconv.Itoa(id)})
}

// GetUserProfile godoc
// @Summary Get a user or trainer profile by ID
// @Description Get a user or trainer profile by ID
// @Tags user
// @Produce json
// @Param id path int true "User or Trainer ID"
// @Success 200 {object} ResponseSuccess{data=domain.User}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/user/{id} [get]
func GetUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid user or trainer ID: " + err.Error()})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, ResponseSuccess{Message: "User found", Data: user})
			return
		}
	}

	for _, trainer := range trainers {
		if trainer.ID == id {
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Trainer found", Data: trainer})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "User or trainer not found with ID " + strconv.Itoa(id)})
}

// UpdateTraining godoc
// @Summary Update a training session by ID
// @Description Update a training session by ID (only for trainers)
// @Tags training
// @Accept json
// @Produce json
// @Param id path int true "Training ID"
// @Param training body domain.Training true "Updated training data"
// @Success 200 {object} ResponseSuccess{data=domain.Training}
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training/{id} [put]
func UpdateTraining(c *gin.Context) {
	trainerID := c.MustGet("user_id").(float64)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid training ID: " + err.Error()})
		return
	}

	var updatedTraining domain.Training
	if err := c.ShouldBindJSON(&updatedTraining); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data: " + err.Error()})
		return
	}

	for i, training := range trainings {
		if training.ID == id {
			if training.TrainerID != int(trainerID) {
				c.JSON(http.StatusForbidden, ResponseError{Error: "Not allowed to update this training"})
				return
			}
			updatedTraining.ID = training.ID
			updatedTraining.TrainerID = training.TrainerID
			trainings[i] = updatedTraining
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Training updated successfully", Data: updatedTraining})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "Training not found with ID " + strconv.Itoa(id)})
}

// DeleteTraining godoc
// @Summary Delete a training session by ID
// @Description Delete a training session by ID (only for trainers)
// @Tags training
// @Produce json
// @Param id path int true "Training ID"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training/{id} [delete]
func DeleteTraining(c *gin.Context) {
	trainerID := c.MustGet("user_id").(float64)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid training ID: " + err.Error()})
		return
	}

	for i, training := range trainings {
		if training.ID == id {
			if training.TrainerID != int(trainerID) {
				c.JSON(http.StatusForbidden, ResponseError{Error: "Not allowed to delete this training"})
				return
			}
			trainings = append(trainings[:i], trainings[i+1:]...)
			for j, trainer := range trainers {
				if trainer.ID == int(trainerID) {
					for k, tid := range trainer.Trainings {
						if tid == id {
							trainers[j].Trainings = append(trainers[j].Trainings[:k], trainers[j].Trainings[k+1:]...)
							break
						}
					}
					break
				}
			}
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Training deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "Training not found with ID " + strconv.Itoa(id)})
}

// UpdateUserProfile godoc
// @Summary Update a user or trainer profile by ID
// @Description Update a user or trainer profile by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User or Trainer ID"
// @Param user body domain.User true "Updated user data"
// @Success 200 {object} ResponseSuccess{data=domain.User}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/user/{id} [put]
func UpdateUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid user or trainer ID: " + err.Error()})
		return
	}

	var updatedUser domain.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data: " + err.Error()})
		return
	}

	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = user.ID
			users[i] = updatedUser
			c.JSON(http.StatusOK, ResponseSuccess{Message: "User profile updated successfully", Data: updatedUser})
			return
		}
	}

	var updatedTrainer domain.Trainer
	if err := c.ShouldBindJSON(&updatedTrainer); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid data: " + err.Error()})
		return
	}

	for i, trainer := range trainers {
		if trainer.ID == id {
			updatedTrainer.ID = trainer.ID
			trainers[i] = updatedTrainer
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Trainer profile updated successfully", Data: updatedTrainer})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "User or trainer not found with ID " + strconv.Itoa(id)})
}

// DeleteUserProfile godoc
// @Summary Delete a user or trainer profile by ID
// @Description Delete a user or trainer profile by ID
// @Tags user
// @Produce json
// @Param id path int true "User or Trainer ID"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/user/{id} [delete]
func DeleteUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid user or trainer ID: " + err.Error()})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, ResponseSuccess{Message: "User profile deleted successfully"})
			return
		}
	}

	for i, trainer := range trainers {
		if trainer.ID == id {
			trainers = append(trainers[:i], trainers[i+1:]...)
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Trainer profile deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "User or trainer not found with ID " + strconv.Itoa(id)})
}

// GetUserSchedule godoc
// @Summary Get the schedule for the current user
// @Description Get the schedule for the current user
// @Tags user
// @Produce json
// @Success 200 {object} ResponseSuccess{data=[]domain.Training}
// @Security BearerAuth
// @Router /protected/user/schedule [get]
func GetUserSchedule(c *gin.Context) {
	userID := c.MustGet("user_id").(float64)
	var userTrainings []domain.Training

	for _, training := range trainings {
		for _, user := range training.Users {
			if user == int(userID) {
				userTrainings = append(userTrainings, training)
				break
			}
		}
	}

	sort.Slice(userTrainings, func(i, j int) bool {
		return userTrainings[i].StartTime.Before(userTrainings[j].StartTime)
	})

	c.JSON(http.StatusOK, ResponseSuccess{Message: "User schedule retrieved", Data: userTrainings})
}

// GetTrainerSchedule godoc
// @Summary Get the schedule for the current trainer
// @Description Get the schedule for the current trainer
// @Tags trainer
// @Produce json
// @Success 200 {object} ResponseSuccess{data=[]domain.Training}
// @Security BearerAuth
// @Router /protected/trainer/schedule [get]
func GetTrainerSchedule(c *gin.Context) {
	trainerID := c.MustGet("user_id").(float64)
	var trainerTrainings []domain.Training

	for _, training := range trainings {
		if training.TrainerID == int(trainerID) {
			trainerTrainings = append(trainerTrainings, training)
		}
	}

	sort.Slice(trainerTrainings, func(i, j int) bool {
		return trainerTrainings[i].StartTime.Before(trainerTrainings[j].StartTime)
	})

	c.JSON(http.StatusOK, ResponseSuccess{Message: "Trainer schedule retrieved", Data: trainerTrainings})
}

// GetAllTrainings godoc
// @Summary Get all available trainings
// @Description Get all available trainings
// @Tags training
// @Produce json
// @Success 200 {object} ResponseSuccess{data=[]domain.Training}
// @Router /trainings [get]
func GetAllTrainings(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseSuccess{Message: "All trainings retrieved", Data: trainings})
}

// GetUsersByTrainingID godoc
// @Summary Get all users by training ID
// @Description Get all users by training ID
// @Tags training
// @Produce json
// @Param training_id path int true "Training ID"
// @Success 200 {object} ResponseSuccess{data=[]domain.User}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Security BearerAuth
// @Router /protected/training/{training_id}/users [get]
func GetUsersByTrainingID(c *gin.Context) {
	trainingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid training ID: " + err.Error()})
		return
	}

	var trainingUsers []domain.User
	for _, training := range trainings {
		if training.ID == trainingID {
			for _, userID := range training.Users {
				for _, user := range users {
					if user.ID == userID {
						trainingUsers = append(trainingUsers, user)
					}
				}
			}
			c.JSON(http.StatusOK, ResponseSuccess{Message: "Users found for training", Data: trainingUsers})
			return
		}
	}

	c.JSON(http.StatusNotFound, ResponseError{Error: "Training not found with ID " + strconv.Itoa(trainingID)})
}
