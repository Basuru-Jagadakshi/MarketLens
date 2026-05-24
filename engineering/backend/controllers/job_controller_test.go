package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"marketlens-go-backend/controllers"
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
)


func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupControllerTestEnv(t *testing.T) (*gorm.DB, *gin.Engine) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=private"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database environment: %v", err)
	}

	t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })

	_ = db.AutoMigrate(&models.JobPost{}, &models.JobMetaData{}, &models.JobType{}, &models.Skill{}, &models.GeoData{})

	db.Create(&models.GeoData{Province: "Western", Latitude: 6.92, Longitude: 79.86})

	repo := repositories.NewJobRepository(db)
	ctrl := controllers.NewJobController(repo)

	r := gin.Default()
	r.POST("/api/jobs", ctrl.CreateJobHandler)
	r.GET("/api/jobs", ctrl.GetAllJobsHandler)
	r.PUT("/api/jobs/:id", ctrl.UpdateJobHandler)
	r.DELETE("/api/jobs/:id", ctrl.DeleteJobHandler)

	return db, r
}


func TestCreateJobHandler_Success(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	input := models.JobPost{
		Employer: "MarketLens Inc",
		JobRole:  "Go Developer",
		MetaData: models.JobMetaData{
			Geo: &models.GeoData{Province: "Western"},
		},
	}
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.JobPost
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotZero(t, response.ID)
	assert.Equal(t, "MarketLens Inc", response.Employer)
}

func TestCreateJobHandler_InvalidJSON(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBufferString("{invalid-json-structure}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload")
}

func TestCreateJobHandler_RepositoryError(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	input := models.JobPost{
		Employer: "MarketLens Inc",
		JobRole:  "Go Developer",
		MetaData: models.JobMetaData{
			Geo: &models.GeoData{Province: "Southern"}, 
		},
	}
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to process job publication data")
}


func TestGetAllJobsHandler_Success(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	db.Create(&models.JobPost{Employer: "Company A", JobRole: "Lead"})
	db.Create(&models.JobPost{Employer: "Company B", JobRole: "Senior"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/jobs", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []models.JobPost
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 2)
}


func TestUpdateJobHandler_Success(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	existingJob := models.JobPost{Employer: "Initial Inc", JobRole: "Junior"}
	db.Create(&existingJob)

	updateInput := models.JobPost{Employer: "Updated Inc", JobRole: "Mid-Level"}
	body, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/api/jobs/%d", existingJob.ID)
	req, _ := http.NewRequest("PUT", path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.JobPost
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Updated Inc", response.Employer)
	assert.Equal(t, "Mid-Level", response.JobRole)
}

func TestUpdateJobHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	
	req, _ := http.NewRequest("PUT", "/api/jobs/abc", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid job ID format parameter")
}

func TestUpdateJobHandler_InvalidJSONPayload(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/jobs/1", bytes.NewBufferString("{broken-json-syntax}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload configuration")
}

func TestUpdateJobHandler_NotFound(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	updateInput := models.JobPost{Employer: "Test", JobRole: "Test"}
	body, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	
	req, _ := http.NewRequest("PUT", "/api/jobs/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to execute modifications on targeted job profile")
}


func TestDeleteJobHandler_Success(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	existingJob := models.JobPost{Employer: "Temp Inc", JobRole: "Contractor"}
	db.Create(&existingJob)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/api/jobs/%d", existingJob.ID)
	req, _ := http.NewRequest("DELETE", path, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedMsg := fmt.Sprintf("Job post with ID %d has been successfully deleted", existingJob.ID)
	assert.Contains(t, w.Body.String(), expectedMsg)
}

func TestDeleteJobHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/jobs/xyz", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid job ID format parameter")
}

func TestDeleteJobHandler_NotFound(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	
	req, _ := http.NewRequest("DELETE", "/api/jobs/999", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to execute deletion on targeted job profile")
}