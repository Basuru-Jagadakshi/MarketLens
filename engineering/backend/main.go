package main

import (
	"marketlens-go-backend/config"
	"marketlens-go-backend/controllers"
	"marketlens-go-backend/repositories"

	"github.com/gin-gonic/gin"
)


func main() { 
	
	config.ConnectDatabase()

	r := gin.Default()

	jobRepo := repositories.NewJobRepository(config.DB)
	jobCtrl := controllers.NewJobController(jobRepo)

	r.POST("/api/jobs", jobCtrl.CreateJobHandler)

	r.Run(":8080")
}
