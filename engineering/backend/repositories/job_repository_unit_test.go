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

func TestDeleteJob_NotFound(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.DeleteJob(404) 

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}