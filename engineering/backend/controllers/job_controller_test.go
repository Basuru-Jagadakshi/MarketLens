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

	_ = db.AutoMigrate(
		&models.JobPost{}, &models.JobMetaData{}, &models.JobType{}, &models.Skill{}, &models.GeoData{},
		&models.MajorGroup{}, &models.SubMajorGroup{}, &models.MinorGroup{}, &models.UnitGroup{}, &models.OccupationGroup{},
		&models.IndustrySector{}, &models.IndustryDivision{}, &models.IndustryGroup{}, &models.IndustryClass{}, &models.IndustrySubclass{},
	)

	db.Create(&models.GeoData{Province: "Western", Latitude: 6.92, Longitude: 79.86})

	repo := repositories.NewJobRepository(db)
	ctrl := controllers.NewJobController(repo)

	r := gin.Default()
	r.DELETE("/api/v1/jobs/:id", ctrl.DeleteJobHandler)
	r.GET("/api/v1/major-groups", ctrl.GetAllMajorGroupsHandler)
	r.GET("/api/v1/major-groups/:id/sub-major-groups", ctrl.GetSubMajorGroupsByMajorGroupHandler)
	r.GET("/api/v1/sub-major-groups/:id/minor-groups", ctrl.GetMinorGroupsBySubMajorGroupHandler)
	r.GET("/api/v1/minor-groups/:id/unit-groups", ctrl.GetUnitGroupsByMinorGroupHandler)
	r.GET("/api/v1/unit-groups/:id/occupation-groups", ctrl.GetOccupationGroupsByUnitGroupHandler)
	r.GET("/api/v1/industry-sectors", ctrl.GetAllIndustrySectorsHandler)
	r.GET("/api/v1/industry-sectors/:id/industry-divisions", ctrl.GetIndustryDivisionsByIndustrySectorHandler)
	r.GET("/api/v1/industry-divisions/:id/industry-groups", ctrl.GetIndustryGroupsByIndustryDivisionHandler)
	r.GET("/api/v1/industry-groups/:id/industry-classes", ctrl.GetIndustryClassesByIndustryGroupHandler)
	r.GET("/api/v1/industry-classes/:id/industry-subclasses", ctrl.GetIndustrySubclassesByIndustryClassHandler)

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


func TestGetAllMajorGroupsHandler_NoDatesReturnsCurrentOnly(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	db.Create(&models.MajorGroup{Name: "Managers", Code: "1"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Managers")
}

func TestGetAllMajorGroupsHandler_ValidDateRangeReturnsResults(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	db.Create(&models.MajorGroup{Name: "Professionals", Code: "2"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups?from-date=2020-01-01&to-date=2030-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Professionals")
}

func TestGetAllMajorGroupsHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups?from-date=01-01-2020&to-date=2030-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetAllMajorGroupsHandler_InvalidToDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups?from-date=2020-01-01&to-date=not-a-date", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid to-date format")
}

func TestGetAllMajorGroupsHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups?from-date=2030-01-01&to-date=2020-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetAllMajorGroupsHandler_OnlyFromDateProvidedFallsBackToCurrentOnly(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	db.Create(&models.MajorGroup{Name: "Solo From Date", Code: "3"})

	w := httptest.NewRecorder()
	// Only from-date supplied - per the handler's `if fromDateStr != "" && toDateStr != ""`
	// check, this should NOT trigger the date-bounded path; it falls back to
	// the plain current-only call, same as if neither were supplied.
	req, _ := http.NewRequest("GET", "/api/v1/major-groups?from-date=2020-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Solo From Date")
}


func TestGetSubMajorGroupsByMajorGroupHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups/abc/sub-major-groups", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid major group id parameter")
}

func TestGetSubMajorGroupsByMajorGroupHandler_NoDatesReturnsCurrentOnly(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	db.Create(&models.SubMajorGroup{MajorGroupID: mg.ID, Name: "Child SMG", Code: "11"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups/1/sub-major-groups", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Child SMG")
}

func TestGetSubMajorGroupsByMajorGroupHandler_ValidDateRange(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	db.Create(&models.SubMajorGroup{MajorGroupID: mg.ID, Name: "Ranged SMG", Code: "11"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups/1/sub-major-groups?from-date=2020-01-01&to-date=2030-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Ranged SMG")
}

func TestGetSubMajorGroupsByMajorGroupHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/major-groups/1/sub-major-groups?from-date=2030-01-01&to-date=2020-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}


func TestGetAllIndustrySectorsHandler_NoDatesReturnsCurrentOnly(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	db.Create(&models.IndustrySector{Name: "Manufacturing", Code: "C"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry-sectors", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Manufacturing")
}

func TestGetAllIndustrySectorsHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry-sectors?from-date=bad&to-date=2030-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}


func TestGetIndustryDivisionsByIndustrySectorHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry-sectors/abc/industry-divisions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid industry sector id parameter")
}

func TestGetIndustryDivisionsByIndustrySectorHandler_NoDatesReturnsCurrentOnly(t *testing.T) {
	db, r := setupControllerTestEnv(t)
	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	db.Create(&models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Child Division", Code: "11"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry-sectors/1/industry-divisions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Child Division")
}

func TestGetTotalVacancyCountHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetTotalVacancyCountHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetTotalVacancyCountHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetTotalVacancyCountHandler_InvalidToDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date=2026-05-01&to-date=not-a-date", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid to-date format")
}

func TestGetTotalVacancyCountHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetTotalVacancyCountHandler_ValidRangeReturnsSum(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 7}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"total_vacancies":10`)
}

func TestGetTotalVacancyCountHandler_NoMatchingJobsReturnsZero(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vacancy-total?from-date=2020-01-01&to-date=2020-12-31", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"total_vacancies":0`)
}


func TestGetOccupationJobCountByDateRangeHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetOccupationJobCountByDateRangeHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetOccupationJobCountByDateRangeHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetOccupationJobCountByDateRangeHandler_InvalidToDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date=2026-05-01&to-date=not-a-date", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid to-date format")
}

func TestGetOccupationJobCountByDateRangeHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetOccupationJobCountByDateRangeHandler_ValidRangeReturnsGroupedResults(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "Professionals", Code: "2"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "21"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "211"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "2111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "21111"}
	db.Create(&og)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Software Engineer", NoOfVacancies: 5}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Professionals")
	assert.Contains(t, w.Body.String(), `"open_job_count":5`)
}

func TestGetOccupationJobCountByDateRangeHandler_NoMajorGroupsReturnsEmptyResults(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupations/by-date-range?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"count":0`)
}

// ---------------------------------------------------------------------------
// GetIndustryJobCountByDateRangeHandler - both dates REQUIRED, no fallback
// ---------------------------------------------------------------------------

func TestGetIndustryJobCountByDateRangeHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetIndustryJobCountByDateRangeHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetIndustryJobCountByDateRangeHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetIndustryJobCountByDateRangeHandler_InvalidToDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date=2026-05-01&to-date=not-a-date", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid to-date format")
}

func TestGetIndustryJobCountByDateRangeHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetIndustryJobCountByDateRangeHandler_ValidRangeReturnsGroupedResults(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Manufacturing", Code: "C"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "C1"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "C11"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "C111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "C1111"}
	db.Create(&subclass)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Factory Worker", NoOfVacancies: 8}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Manufacturing")
	assert.Contains(t, w.Body.String(), `"open_job_count":8`)
}

func TestGetIndustryJobCountByDateRangeHandler_NoIndustrySectorsReturnsEmptyResults(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industries/by-date-range?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"count":0`)
}


func TestGetEmploymentSectorByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/employment-sector?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetEmploymentSectorByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/employment-sector?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetEmploymentSectorByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/employment-sector?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetEmploymentSectorByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/employment-sector?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetEmploymentSectorByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/employment-sector?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetEmploymentSectorByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/employment-sector?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetEmploymentSectorByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/employment-sector?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetEmploymentSectorByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	sector := models.EmploymentSector{Sector: "Private"}
	db.Create(&sector)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EmploymentSectorID: sector.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/employment-sector?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Private")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetEmploymentSectorByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	empSector := models.EmploymentSector{Sector: "Government"}
	db.Create(&empSector)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, EmploymentSectorID: empSector.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/employment-sector?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Government")
	assert.Contains(t, w.Body.String(), `"open_job_count":4`)
}


func TestGetExperienceByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/experience?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetExperienceByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/experience?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetExperienceByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/experience?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetExperienceByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/experience?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetExperienceByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/experience?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetExperienceByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/experience?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetExperienceByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/experience?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetExperienceByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	exp := models.Experience{Name: "Entry Level"}
	db.Create(&exp)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, ExperienceID: exp.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/experience?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Entry Level")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetExperienceByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	exp := models.Experience{Name: "Senior"}
	db.Create(&exp)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, ExperienceID: exp.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/experience?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Senior")
	assert.Contains(t, w.Body.String(), `"open_job_count":4`)
}

// ---------------------------------------------------------------------------
// GetProvinceByLevelHandler - standard/level/id path params +
// required from-date/to-date
// ---------------------------------------------------------------------------

func TestGetProvinceByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/province?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetProvinceByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/province?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetProvinceByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/province?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetProvinceByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/province?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetProvinceByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/province?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetProvinceByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/province?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetProvinceByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/province?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetProvinceByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	// Note: setupControllerTestEnv already seeds a "Western" GeoData row -
	// reuse it here rather than creating a duplicate.
	var western models.GeoData
	db.Where("province = ?", "Western").First(&western)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 5}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GeoDataID: western.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/province?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Western")
	assert.Contains(t, w.Body.String(), `"open_job_count":5`)
}

func TestGetProvinceByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	var central models.GeoData
	db.Where("province = ?", "Central").First(&central)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 3}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, GeoDataID: central.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/province?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Central")
	assert.Contains(t, w.Body.String(), `"open_job_count":3`)
}


func TestGetEducationLevelByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetEducationLevelByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetEducationLevelByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetEducationLevelByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/education?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetEducationLevelByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/education?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetEducationLevelByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/education?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetEducationLevelByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/education?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetEducationLevelByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	edu := models.EducationLevel{Level: "Degree"}
	db.Create(&edu)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EducationLevelID: edu.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/education?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Degree")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetEducationLevelByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	edu := models.EducationLevel{Level: "A/L"}
	db.Create(&edu)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 2}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, EducationLevelID: edu.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/education?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "A/L")
	assert.Contains(t, w.Body.String(), `"open_job_count":2`)
}

func TestGetFormalityByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/formality?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetFormalityByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/formality?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetFormalityByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/formality?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetFormalityByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/formality?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetFormalityByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/formality?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetFormalityByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/formality?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetFormalityByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/formality?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetFormalityByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	formality := models.Formality{FormalityType: "Formal"}
	db.Create(&formality)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, FormalityID: formality.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/formality?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Formal")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetFormalityByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	formality := models.Formality{FormalityType: "Informal"}
	db.Create(&formality)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 2}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, FormalityID: formality.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/formality?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Informal")
	assert.Contains(t, w.Body.String(), `"open_job_count":2`)
}

// ---------------------------------------------------------------------------
// GetGenderByLevelHandler - standard/level/id path params +
// required from-date/to-date
// ---------------------------------------------------------------------------

func TestGetGenderByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/gender?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetGenderByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/gender?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetGenderByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/gender?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetGenderByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/gender?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetGenderByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/gender?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetGenderByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/gender?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetGenderByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/gender?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetGenderByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	gender := models.Gender{GenderType: "Male"}
	db.Create(&gender)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GenderID: gender.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/gender?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Male")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetGenderByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	gender := models.Gender{GenderType: "Female"}
	db.Create(&gender)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 2}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, GenderID: gender.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/gender?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Female")
	assert.Contains(t, w.Body.String(), `"open_job_count":2`)
}

func TestGetVocationalEducationByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/vocational-education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetVocationalEducationByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/vocational-education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetVocationalEducationByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/vocational-education?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetVocationalEducationByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/vocational-education?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetVocationalEducationByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/vocational-education?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetVocationalEducationByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/vocational-education?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetVocationalEducationByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/vocational-education?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetVocationalEducationByLevelHandler_ValidRequestReturnsBreakdown(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	nvq := models.VocationalEducation{Level: "NVQ 4"}
	db.Create(&nvq)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, VocationalEducationID: nvq.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/vocational-education?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "NVQ 4")
	assert.Contains(t, w.Body.String(), `"open_job_count":6`)
}

func TestGetVocationalEducationByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	nvq := models.VocationalEducation{Level: "NVQ 5"}
	db.Create(&nvq)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 2}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, VocationalEducationID: nvq.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/vocational-education?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "NVQ 5")
	assert.Contains(t, w.Body.String(), `"open_job_count":2`)
}


func TestGetRemoteOnSiteByLevelHandler_InvalidStandard(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/not-a-real-standard/major-group/1/remote-onsite?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid standard")
}

func TestGetRemoteOnSiteByLevelHandler_InvalidLevelForOccupation(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/remote-onsite?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for standard 'occupation'")
}

func TestGetRemoteOnSiteByLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/remote-onsite?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetRemoteOnSiteByLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetRemoteOnSiteByLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetRemoteOnSiteByLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetRemoteOnSiteByLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetRemoteOnSiteByLevelHandler_ValidRequestReturnsSplitCounts(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	remoteJob := models.JobPost{JobTypeID: jobType.ID, JobRole: "Remote Role", NoOfVacancies: 4, IsRemote: true}
	db.Create(&remoteJob)
	db.Create(&models.JobMetaData{JobPostID: remoteJob.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	onSiteJob := models.JobPost{JobTypeID: jobType.ID, JobRole: "On-Site Role", NoOfVacancies: 2, IsRemote: false}
	db.Create(&onSiteJob)
	db.Create(&models.JobMetaData{JobPostID: onSiteJob.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"remote_count":4`)
	assert.Contains(t, w.Body.String(), `"on_site_count":2`)
}

func TestGetRemoteOnSiteByLevelHandler_NoMatchingJobsReturnsZeroForBoth(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/remote-onsite?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"remote_count":0`)
	assert.Contains(t, w.Body.String(), `"on_site_count":0`)
}

func TestGetRemoteOnSiteByLevelHandler_WorksForIndustryStandard(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	sector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&sector)
	division := models.IndustryDivision{IndustrySectorID: sector.ID, Name: "Division", Code: "11"}
	db.Create(&division)
	group := models.IndustryGroup{IndustryDivisionID: division.ID, Name: "Group", Code: "111"}
	db.Create(&group)
	class := models.IndustryClass{IndustryGroupID: group.ID, Name: "Class", Code: "1111"}
	db.Create(&class)
	subclass := models.IndustrySubclass{IndustryClassID: class.ID, Name: "Subclass", Code: "11111"}
	db.Create(&subclass)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 3, IsRemote: true}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: subclass.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/industry/industry-sector/1/remote-onsite?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"remote_count":3`)
	assert.Contains(t, w.Body.String(), `"on_site_count":0`)
}

// ---------------------------------------------------------------------------
// GetTop15SkillsByOccupationLevelHandler - occupation-only, level/id path
// params + required from-date/to-date
// ---------------------------------------------------------------------------

func TestGetTop15SkillsByOccupationLevelHandler_InvalidLevel(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/top-15-skills?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for occupation")
}

func TestGetTop15SkillsByOccupationLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/top-15-skills?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetTop15SkillsByOccupationLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetTop15SkillsByOccupationLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetTop15SkillsByOccupationLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetTop15SkillsByOccupationLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetTop15SkillsByOccupationLevelHandler_ValidRequestReturnsSkills(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	skill := models.Skill{Skill: "Go"}
	db.Create(&skill)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 5}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go")
	assert.Contains(t, w.Body.String(), `"open_job_count":5`)
}

func TestGetTop15SkillsByOccupationLevelHandler_NoMatchingSkillsReturnsEmptyResults(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-15-skills?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"count":0`)
}

// ---------------------------------------------------------------------------
// GetTopHiringEmployersByOccupationLevelHandler - occupation-only, level/id
// path params + required from-date/to-date
// ---------------------------------------------------------------------------

func TestGetTopHiringEmployersByOccupationLevelHandler_InvalidLevel(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/not-a-real-level/1/top-hiring-employers?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level for occupation")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_InvalidIDFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/abc/top-hiring-employers?from-date=2026-05-01&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id parameter")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_MissingFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "from-date query parameter is required")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_MissingToDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?from-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date query parameter is required")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_InvalidFromDateFormat(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?from-date=01-05-2026&to-date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from-date format")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_ToDateBeforeFromDate(t *testing.T) {
	_, r := setupControllerTestEnv(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?from-date=2026-08-01&to-date=2026-05-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to-date must not be before from-date")
}

func TestGetTopHiringEmployersByOccupationLevelHandler_ValidRequestReturnsEmployers(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	emp := models.Employer{Name: "Acme Corp"}
	db.Create(&emp)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{EmployerID: emp.ID, JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 9}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Acme Corp")
	assert.Contains(t, w.Body.String(), `"open_job_count":9`)
}

func TestGetTopHiringEmployersByOccupationLevelHandler_NoMatchingJobsReturnsEmptyResults(t *testing.T) {
	db, r := setupControllerTestEnv(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)

	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/occupation/major-group/1/top-hiring-employers?from-date="+fromDate+"&to-date="+toDate, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"count":0`)
}