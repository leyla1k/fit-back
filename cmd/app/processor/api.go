package processor

import (
	"github.com/folklinoff/fitness-app/internal/handler"
	middleware "github.com/folklinoff/fitness-app/internal/middleware/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func api() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.POST("/login/:user_type", handler.Login)
	r.POST("/register/:user_type", handler.Register)
	r.GET("/trainings", handler.GetAllTrainings)

	// Protected routes
	protected := r.Group("/protected")
	protected.Use(middleware.AuthenticationMiddleware())
	{
		protected.GET("/profile", handler.Profile)
		protected.POST("/training", handler.CreateTraining)
		protected.POST("/training/:id/register", handler.RegisterUserForTraining)
		protected.GET("/training/:id", handler.GetTrainingByID)
		protected.PUT("/training/:id", handler.UpdateTraining)
		protected.DELETE("/training/:id", handler.DeleteTraining)
		protected.GET("/user/:id", handler.GetUserProfile)
		protected.PUT("/user/:id", handler.UpdateUserProfile)
		protected.DELETE("/user/:id", handler.DeleteUserProfile)
		protected.GET("/user/schedule", handler.GetUserSchedule)
		protected.GET("/trainer/schedule", handler.GetTrainerSchedule)
		protected.GET("/training/:id/users", handler.GetUsersByTrainingID)
	}

	return r
}
