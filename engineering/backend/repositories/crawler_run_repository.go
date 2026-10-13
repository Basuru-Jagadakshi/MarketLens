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
		// can we move these string values up in one place. so we can access through variables and can be easily changes from one place later.
		Status: "RUNNING",
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// starting a run marks any RUNNING run as FAILED, but that run's crawler process keeps going.
		// it later calls reconcile with its own id, which ends every open job with crawler_run_id <> oldRun,
		// including jobs the new run just confirmed as live, and then marks itself COMPLETED.
		// so the old run goes RUNNING -> FAILED -> COMPLETED and live vacancies get closed.
		// this happens on any overlap (manual + scheduled run, a slow run overlapping the next schedule, retries).
		//
		// failing old RUNNING runs is probably meant to clean up crashed runs that would otherwise block new
		// runs forever. suggest separating the two cases:
		//   - RUNNING and recent -> active: refuse the new start (409)
		//   - RUNNING and older than the max expected run time -> crashed: mark FAILED, then start
		// also:
		//   - reconcile and complete should refuse runs that aren't RUNNING (safety net)
		//   - a partial unique index (WHERE status = 'RUNNING') prevents two concurrent starts both succeeding
		// this enforces "one crawler run at a time". supporting concurrent runs would also need reconcile
		// to be scoped per run/source, so one-at-a-time is the simpler option. open to discuss.

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

// if this supports for any state change why the function name is "CopmleteCrawlerRun"?
func (r *JobRepository) CompleteCrawlerRun(id uint, status string) error {
	now := time.Now()
	return r.db.Model(&models.CrawlerRun{}).Where("id = ?", id).Updates(map[string]interface{}{
		"finished_at": &now,
		"status":      status, // 'COMPLETED' or 'FAILED'
	}).Error
}
