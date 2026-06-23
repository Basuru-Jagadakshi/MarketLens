package controllers_test

import (
	"bytes"
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
	r.DELETE("/api/v1/jobs/:id", ctrl.DeleteJobHandler)

	return db, r
}

func TestCreateJobHandler_InvalidJSON(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/jobs", bytes.NewBufferString("{invalid-json-structure}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload")
}

func TestUpdateJobHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	
	req, _ := http.NewRequest("PUT", "/api/v1/jobs/abc", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid job ID format parameter")
}

func TestUpdateJobHandler_InvalidJSONPayload(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/jobs/1", bytes.NewBufferString("{broken-json-syntax}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload configuration")
}

func TestDeleteJobHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/jobs/xyz", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid job ID format parameter")
}

func TestDeleteJobHandler_NotFound(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	
	req, _ := http.NewRequest("DELETE", "/api/v1/jobs/999", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to execute deletion on targeted job profile")
}