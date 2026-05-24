//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"marketlens-go-backend/controllers"
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
)


func setupIntegrationTestStack(t *testing.T) *gin.Engine {
	
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&parseTime=True"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to spin up integration test data layer: %v", err)
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&models.JobMetaData{}); err == nil {
		if field, ok := stmt.Schema.FieldsByDBName["posted_at"]; ok {
			field.DataType = "time"
			field.GORMDataType = "time"
		}
	}

	_ = db.AutoMigrate(&models.JobPost{}, &models.JobMetaData{}, &models.JobType{}, &models.Skill{}, &models.GeoData{})

	db.Create(&models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612})
	db.Create(&models.GeoData{Province: "Central", Latitude: 7.2906, Longitude: 80.6337})

	repo := repositories.NewJobRepository(db)
	ctrl := controllers.NewJobController(repo)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/api/jobs", ctrl.CreateJobHandler)
	r.GET("/api/jobs", ctrl.GetAllJobsHandler)
	r.PUT("/api/jobs/:id", ctrl.UpdateJobHandler)
	r.DELETE("/api/jobs/:id", ctrl.DeleteJobHandler)

	return r
}

func TestJobLifecycle_Integration(t *testing.T) {
	r := setupIntegrationTestStack(t)

	newJobPayload := models.JobPost{
		Employer: "MarketLens International",
		JobRole:  "Senior Backend Engineer (Go)",
		JobType:  models.JobType{Name: "Full-Time"},
		MetaData: models.JobMetaData{
			Source: "Internal Portal",
			Geo:    &models.GeoData{Province: "Western"}, 
		},
		Skills: []models.Skill{
			{Name: "Golang"},
			{Name: "Docker"},
		},
	}
	createBody, _ := json.Marshal(newJobPayload)

	wCreate := httptest.NewRecorder()
	reqCreate, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(createBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, wCreate.Code)
	
	var createdJob models.JobPost
	json.Unmarshal(wCreate.Body.Bytes(), &createdJob)
	assert.NotZero(t, createdJob.ID, "The database layer should have committed a valid auto-increment ID index")
	assert.Equal(t, "Senior Backend Engineer (Go)", createdJob.JobRole)
	assert.Len(t, createdJob.Skills, 2)

	updateJobPayload := models.JobPost{
		Employer: "MarketLens International",
		JobRole:  "Lead Distributed Systems Architect", 
		JobType:  models.JobType{Name: "Full-Time"},
	}
	updateBody, _ := json.Marshal(updateJobPayload)

	wUpdate := httptest.NewRecorder()
	updatePath := fmt.Sprintf("/api/jobs/%d", createdJob.ID)
	reqUpdate, _ := http.NewRequest("PUT", updatePath, bytes.NewBuffer(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)
	var updatedJob models.JobPost
	json.Unmarshal(wUpdate.Body.Bytes(), &updatedJob)
	assert.Equal(t, "Lead Distributed Systems Architect", updatedJob.JobRole)

	wGetAll := httptest.NewRecorder()
	reqGetAll, _ := http.NewRequest("GET", "/api/jobs", nil)
	r.ServeHTTP(wGetAll, reqGetAll)

	assert.Equal(t, http.StatusOK, wGetAll.Code)
	var completeList []models.JobPost
	json.Unmarshal(wGetAll.Body.Bytes(), &completeList)
	assert.True(t, len(completeList) >= 1)

	wDelete := httptest.NewRecorder()
	deletePath := fmt.Sprintf("/api/jobs/%d", createdJob.ID)
	reqDelete, _ := http.NewRequest("DELETE", deletePath, nil)
	r.ServeHTTP(wDelete, reqDelete)

	assert.Equal(t, http.StatusOK, wDelete.Code)
	assert.Contains(t, wDelete.Body.String(), "successfully deleted")
}