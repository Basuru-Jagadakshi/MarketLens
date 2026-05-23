package repositories_test

import (
	"errors"
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



func setupTestDB(t *testing.T) (*gorm.DB, *repositories.JobRepository) {
	
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=private"), &gorm.Config{
		NowFunc: func() time.Time {
            return time.Now().UTC()
        },
	})
	if err != nil {
		t.Fatalf("Failed to initialize temporary in-memory database workspace: %v", err)
	}

	t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })


	if err := db.Callback().Create().Before("gorm:create").Register("sqlite_type_fix", func(d *gorm.DB) {
		
	}); err != nil {
		t.Fatalf("Failed to register driver override schema callbacks: %v", err)
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&models.JobMetaData{}); err == nil {
		if field, ok := stmt.Schema.FieldsByDBName["posted_at"]; ok {
			field.DataType = "time"
			field.GORMDataType = "time"
		}
	}

	err = db.AutoMigrate(
		&models.JobPost{},
		&models.JobMetaData{},
		&models.JobType{},
		&models.Skill{},
		&models.GeoData{},
	)
	if err != nil {
		t.Fatalf("Schema generation auto-migration failed: %v", err)
	}

	db.Create(&models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612})
	db.Create(&models.GeoData{Province: "Central", Latitude: 7.2906, Longitude: 80.6337})

	repo := repositories.NewJobRepository(db)
	return db, repo
}


func TestCreateJob_Success(t *testing.T) {
	_, repo := setupTestDB(t)

	input := models.JobPost{
		Employer: "MarketLens Analytics Inc.",
		JobRole:  "Software Engineer",
		JobType: models.JobType{
			Name: "Full Time",
		},
		MetaData: models.JobMetaData{
			Source:              "LinkedIn",
			StandarizedCategory: "Engineering",
			Geo: &models.GeoData{
				Province: "Western", 
			},
		},
		Skills: []models.Skill{
			{Name: "Go"},
			{Name: "GORM"},
		},
	}

	result, err := repo.CreateJob(input)

	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "MarketLens Analytics Inc.", result.Employer)
	assert.Equal(t, "Full Time", result.JobType.Name)
	assert.NotNil(t, result.MetaData.Geo)
	assert.Equal(t, "Western", result.MetaData.Geo.Province)
	assert.Len(t, result.Skills, 2)
}


func TestCreateJob_InvalidProvince_ReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	input := models.JobPost{
		Employer: "MarketLens Systems",
		JobRole:  "Architect",
		MetaData: models.JobMetaData{
			Geo: &models.GeoData{
				Province: "NonExistentProvince", 
			},
		},
	}

	_, err := repo.CreateJob(input)

	assert.Error(t, err)
	assert.Equal(t, "provided province does not match system geo registers", err.Error())
}


func TestCreateJob_FirstOrCreate_Skills_Deduplication(t *testing.T) {
	db, repo := setupTestDB(t)

	db.Create(&models.Skill{Name: "Go"})

	job1 := models.JobPost{Employer: "Company A", JobRole: "Dev", Skills: []models.Skill{{Name: "Go"}}}
	job2 := models.JobPost{Employer: "Company B", JobRole: "Dev", Skills: []models.Skill{{Name: "Go"}}}

	_, err1 := repo.CreateJob(job1)
	_, err2 := repo.CreateJob(job2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	var count int64
	db.Model(&models.Skill{}).Where("name = ?", "Go").Count(&count)
	assert.Equal(t, int64(1), count, "The system should reuse the existing skill record instead of generating clones")
}


func TestGetAllJobs_EmptyAndPopulated(t *testing.T) {
	db, repo := setupTestDB(t)

	emptyResults, err := repo.GetAllJobs()
	assert.NoError(t, err)
	assert.Len(t, emptyResults, 0)

	db.Create(&models.JobPost{Employer: "Team Alpha", JobRole: "Engineer"})
	db.Create(&models.JobPost{Employer: "Team Beta", JobRole: "Designer"})

	populatedResults, err := repo.GetAllJobs()
	assert.NoError(t, err)
	assert.Len(t, populatedResults, 2)
	assert.Equal(t, "Team Alpha", populatedResults[0].Employer)
	assert.Equal(t, "Team Beta", populatedResults[1].Employer)
}


func TestUpdateJob_Success(t *testing.T) {
	_, repo := setupTestDB(t)

	initialJob, _ := repo.CreateJob(models.JobPost{
		Employer: "MarketLens Inc.",
		JobRole:  "Junior Engineer",
		JobType:  models.JobType{Name: "Part Time"},
		Skills:   []models.Skill{{Name: "Python"}},
	})

	updateInput := models.JobPost{
		Employer: "MarketLens Core Analytics",
		JobRole:  "Senior Engineer",
		JobType:  models.JobType{Name: "Full Time"},
		Skills:   []models.Skill{{Name: "Go"}},
	}

	updatedJob, err := repo.UpdateJob(initialJob.ID, updateInput)

	assert.NoError(t, err)
	assert.Equal(t, initialJob.ID, updatedJob.ID)
	assert.Equal(t, "MarketLens Core Analytics", updatedJob.Employer)
	assert.Equal(t, "Senior Engineer", updatedJob.JobRole)
	assert.Equal(t, "Full Time", updatedJob.JobType.Name)

	assert.Len(t, updatedJob.Skills, 2)

	skillNames := []string{updatedJob.Skills[0].Name, updatedJob.Skills[1].Name}
	assert.Contains(t, skillNames, "Python")
	assert.Contains(t, skillNames, "Go")
}


func TestUpdateJob_RecordNotFound(t *testing.T) {
	_, repo := setupTestDB(t)

	missingInput := models.JobPost{Employer: "Ghost Corp", JobRole: "None"}
	_, err := repo.UpdateJob(999, missingInput) 

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}


func TestDeleteJob_Success(t *testing.T) {
	db, repo := setupTestDB(t)

	createdJob, _ := repo.CreateJob(models.JobPost{
		Employer: "Delete Target Corp",
		JobRole:  "Temporary Role",
		Skills:   []models.Skill{{Name: "Ephemeral Skill"}},
	})

	deletedID, err := repo.DeleteJob(createdJob.ID)

	assert.NoError(t, err)
	assert.Equal(t, createdJob.ID, deletedID)

	var searchJob models.JobPost
	err = db.First(&searchJob, createdJob.ID).Error
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	var joinCount int64
	db.Table("job_skills").Where("job_post_id = ?", createdJob.ID).Count(&joinCount)
	assert.Equal(t, int64(0), joinCount, "Junction table references should be deleted completely")
}

func TestDeleteJob_NotFound(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.DeleteJob(404) 

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}