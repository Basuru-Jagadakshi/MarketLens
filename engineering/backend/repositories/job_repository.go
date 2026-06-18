package repositories

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"time"

	"gorm.io/gorm"
)



type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateCrawlerRun() (models.CrawlerRun, error) {
	run := models.CrawlerRun{
		StartedAt: time.Now(),
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

func (r *JobRepository) GetJobsByBucketKeys(bucketKeys []string) ([]models.JobPost, error) {
	var jobs []models.JobPost

	err := r.db.Distinct("job_posts.*").
		Joins("JOIN lsh_index ON lsh_index.job_post_id = job_posts.id").
		Joins("JOIN job_metadata ON job_metadata.job_post_id = job_posts.id").
		Where("lsh_index.bucket_key IN ? AND job_metadata.end_date IS NULL", bucketKeys).
		Preload("MetaData"). 
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// func (r *JobRepository) BatchSaveNewJobs(jobs []models.JobPost, lshIndexRecords []models.LshIndex) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {

// 		if len(jobs) > 0 {
// 			if err := tx.CreateInBatches(&jobs, 100).Error; err != nil {
// 				return err
// 			}
// 		}

// 		for i := range lshIndexRecords {
// 			lshIndexRecords[i].JobPostID = jobs[i].ID
// 		}

// 		if len(lshIndexRecords) > 0 {
// 			if err := tx.CreateInBatches(&lshIndexRecords, 500).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

func (r *JobRepository) BatchSaveNewJobs(jobs []models.JobPost, lshIndexRecords []models.LshIndex) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		
		// Map to store newly generated IDs by their iteration index
		generatedJobIDs := make(map[int]uint)

		// 1. Process and save jobs individual loops using explicit pointer targets
		for i := range jobs {
			job := &jobs[i] 

			// A. Geo/Province Lookup
			if job.MetaData.Geo != nil && job.MetaData.Geo.Province != "" {
				var existingGeo models.GeoData
				err := tx.Where("province = ?", job.MetaData.Geo.Province).First(&existingGeo).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job index %d error: province '%s' not registered", i, job.MetaData.Geo.Province)
				} else if err != nil {
					return err
				}
				job.MetaData.GeoID = &existingGeo.ID
				job.MetaData.Geo = nil 
			}

			// B. JobType Lookup
			if job.JobType.Name != "" {
				var jt models.JobType
				if err := tx.Where(models.JobType{Name: job.JobType.Name}).FirstOrCreate(&jt).Error; err != nil {
					return err
				}
				job.JobTypeID = jt.ID
				job.JobType = models.JobType{} 
			}

			// C. Skills Lookup
			var linkedSkills []models.Skill
			for _, s := range job.Skills {
				if s.Name == "" {
					continue
				}
				var skill models.Skill
				if err := tx.Where(models.Skill{Name: s.Name}).FirstOrCreate(&skill).Error; err != nil {
					return err
				}
				linkedSkills = append(linkedSkills, skill)
			}
			job.Skills = linkedSkills
			
			// D. Save to DB 
			if err := tx.Create(job).Error; err != nil {
				return err
			}

			// Capture the guaranteed new database ID directly from the active pointer context
			generatedJobIDs[i] = job.ID
		}

		// 2. Map generated database primary keys back onto LSH indices array safely
		for idx := range lshIndexRecords {
			jobGroupIndex := idx / 8 
			
			// Assign the true database ID from our tracking map
			if trueID, exists := generatedJobIDs[jobGroupIndex]; exists {
				lshIndexRecords[idx].JobPostID = trueID
			} else {
				return fmt.Errorf("failed to match LSH index %d to a valid generated Job ID", idx)
			}
		}

		// 3. Bulk insert the updated LSH records
		if len(lshIndexRecords) > 0 {
			if err := tx.Omit("JobPost").CreateInBatches(&lshIndexRecords, 500).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *JobRepository) BatchUpdateDuplicateJobs(updates []models.JobMetaData) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, up := range updates {
			err := tx.Table("job_metadata").
				Where("job_post_id = ?", up.JobPostID).
				Update("crawler_run_id", up.CrawlerRunID).Error
				
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *JobRepository) ReconcileStaleVacancies(currentRunID uint) (int64, error) {
	var affectedRows int64
	now := time.Now()

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var staleJobIDs []uint
		err := tx.Model(&models.JobMetaData{}).
			Where("crawler_run_id <> ? AND end_date IS NULL", currentRunID).
			Pluck("job_post_id", &staleJobIDs).Error
		if err != nil {
			return err
		}

		if len(staleJobIDs) == 0 {
			return nil
		}

		if err := tx.Where("job_post_id IN ?", staleJobIDs).Delete(&models.LshIndex{}).Error; err != nil {
			return err
		}

		result := tx.Model(&models.JobMetaData{}).
			Where("job_post_id IN ?", staleJobIDs).
			Updates(map[string]interface{}{
				"end_date": &now,
			})
		if result.Error != nil {
			return result.Error
		}

		affectedRows = result.RowsAffected
		return nil
	})

	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (r *JobRepository) CreateJob(input models.JobPost) (models.JobPost, error) {
    
    var completeJob models.JobPost  // declare outside

    err := r.db.Transaction(func(tx *gorm.DB) error {

        if input.MetaData.Geo != nil && input.MetaData.Geo.Province != "" {
            var existingGeo models.GeoData
            err := tx.Where("province = ?", input.MetaData.Geo.Province).First(&existingGeo).Error
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return errors.New("provided province does not match system geo registers")
            } else if err != nil {
                return err
            }
            input.MetaData.GeoID = &existingGeo.ID
            input.MetaData.Geo = nil
        }

        if input.JobType.Name != "" {
            var jobType models.JobType
            if err := tx.Where(models.JobType{Name: input.JobType.Name}).FirstOrCreate(&jobType).Error; err != nil {
                return err
            }
            input.JobTypeID = jobType.ID
            input.JobType = models.JobType{}
        }

        var linkedSkills []models.Skill
        for _, s := range input.Skills {
            if s.Name == "" {
                continue
            }
            var skill models.Skill
            if err := tx.Where(models.Skill{Name: s.Name}).FirstOrCreate(&skill).Error; err != nil {
                return err
            }
            linkedSkills = append(linkedSkills, skill)
        }
        input.Skills = linkedSkills

        if err := tx.Create(&input).Error; err != nil {
            return err
        }

        // Fetch complete job inside transaction while ID is guaranteed set
        return tx.Preload("JobType").Preload("MetaData.Geo").Preload("Skills").First(&completeJob, input.ID).Error
    })

    if err != nil {
        return models.JobPost{}, err
    }

    return completeJob, nil  // completeJob is fully populated
}


//job get all API
func (r *JobRepository) GetAllJobs() ([]models.JobPost, error) {
	
	var jobs []models.JobPost

	err := r.db.Preload("JobType").
		Preload("MetaData").
		Preload("MetaData.Geo").
		Preload("Skills").
		Joins("JOIN job_metadata ON job_metadata.job_post_id = job_posts.id AND job_metadata.end_date IS NULL").
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}


//job update API
func (r *JobRepository) UpdateJob(id uint, input models.JobPost) (models.JobPost, error) {
	
	var existingJob models.JobPost
	
	//check whether job post exists or not
	if err := r.db.Preload("MetaData").First(&existingJob, id).Error; err != nil {
		return models.JobPost{}, err
	}

	input.ID = id
	input.MetaData.JobPostID = id
	input.MetaData.ID = existingJob.MetaData.ID

	err := r.db.Transaction(func(tx *gorm.DB) error {
		
		//This block checks whether province in json payload matches with the values in DB
		if input.MetaData.Geo != nil && input.MetaData.Geo.Province != "" {
			var existingGeo models.GeoData
			err := tx.Where("province = ?", input.MetaData.Geo.Province).First(&existingGeo).Error
			if err != nil {
				return errors.New("provided province does not match system geo registers")
			}
			input.MetaData.GeoID = &existingGeo.ID
			input.MetaData.Geo = nil
		}

		//This block checks whether job type in json payload matches with the values in DB.
		//If there is no value add it as a new value.
		if input.JobType.Name != "" {
			var jobType models.JobType
			if err := tx.Where(models.JobType{Name: input.JobType.Name}).FirstOrCreate(&jobType).Error; err != nil {
				return err
			}
			input.JobTypeID = jobType.ID
			input.JobType = models.JobType{}
		}

		//This block checks whether skills in json payload matches with the values in DB
		//If there is no value add it as a new value.
		var linkedSkills []models.Skill
		for _, s := range input.Skills {
			if s.Name == "" {
				continue
			}
			var skill models.Skill
			if err := tx.Where(models.Skill{Name: s.Name}).FirstOrCreate(&skill).Error; err != nil {
				return err
			}
			linkedSkills = append(linkedSkills, skill)
		}
		input.Skills = linkedSkills

		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&input).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.JobPost{}, err
	}

	var completeJob models.JobPost
	r.db.Preload("JobType").Preload("MetaData").Preload("MetaData.Geo").Preload("Skills").First(&completeJob, id)

	return completeJob, nil
}


//job update API
func (r *JobRepository) DeleteJob(id uint) (uint, error) {

	var job models.JobPost

	//check whether job post exists or not
	if err := r.db.First(&job, id).Error; err != nil {
		return 0,err
	}

	if err := r.db.Select("Skills").Delete(&job).Error; err != nil {
		return 0,err
	}

	return job.ID, nil
}