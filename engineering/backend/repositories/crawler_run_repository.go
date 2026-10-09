package repositories

import (
	"marketlens-go-backend/models"
	"time"

	"gorm.io/gorm"
)

func (r *JobRepository) CreateCrawlerRun() (models.CrawlerRun, error) {
	now := time.Now()
	run := models.CrawlerRun{
		StartedAt: &now,
		Status:    "RUNNING",
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		err := tx.Model(&models.CrawlerRun{}).
			Where("status = ?", "RUNNING").
			Updates(map[string]interface{}{
				"status":      "FAILED",
				"finished_at": &now,
			}).Error
		if err != nil {
			return err
		}

		if err := tx.Create(&run).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.CrawlerRun{}, err
	}

	return run, nil
}

func (r *JobRepository) CompleteCrawlerRun(id uint, status string) error {
	now := time.Now()
	return r.db.Model(&models.CrawlerRun{}).Where("id = ?", id).Updates(map[string]interface{}{
		"finished_at": &now,
		"status":      status, // 'COMPLETED' or 'FAILED'
	}).Error
}
