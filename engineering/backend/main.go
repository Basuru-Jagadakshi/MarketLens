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

	v1 := r.Group("/api/v1")
	{
		v1.POST("/jobs", jobCtrl.CreateJobHandler)
		v1.GET("/jobs", jobCtrl.GetAllJobsHandler)
		v1.PUT("/jobs/:id", jobCtrl.UpdateJobHandler)
		v1.DELETE("/jobs/:id", jobCtrl.DeleteJobHandler)
	}

	r.Run(":8080")
}
