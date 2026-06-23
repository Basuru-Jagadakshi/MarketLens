package repositories

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"time"
	"math"

	"gorm.io/gorm"
)



type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

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

func (r *JobRepository) GetJobsByBucketKeys(bucketKeys []string) ([]models.JobPost, error) {
	var jobs []models.JobPost

	err := r.db.Distinct("job_post.*").
		Joins("JOIN lsh_index ON lsh_index.job_post_id = job_post.id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("lsh_index.bucket_key IN ? AND meta_data.end_date IS NULL", bucketKeys).
		Preload("Employer").
		Preload("JobType").
		Preload("Skills").
		Preload("MetaData"). 
		Preload("MetaData.AiVersion").
		Preload("MetaData.EducationLevel").
		Preload("MetaData.GeoData").
		Preload("MetaData.Industry").
		Preload("MetaData.Occupation").
		Preload("MetaData.Source").
		Preload("MetaData.Experience").
		Preload("MetaData.CrawlerRun").
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) BatchSaveNewJobs(jobs []models.JobPost, lshIndexRecords []models.LshIndex) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		generatedJobIDs := make(map[int]uint)

		for i := range jobs {
			job := &jobs[i]

			// ── A. Employer (FirstOrCreate) ───────────────────────────────────
			if job.Employer != nil && job.Employer.Name != "" {
				var employer models.Employer
				if err := tx.Where(models.Employer{Name: job.Employer.Name}).
					FirstOrCreate(&employer).Error; err != nil {
					return fmt.Errorf("job[%d] employer lookup failed: %w", i, err)
				}
				job.EmployerID = &employer.ID
				job.Employer = nil
			}

			// ── B. JobType (FirstOrCreate) ────────────────────────────────────
			if job.JobType != nil && job.JobType.Type != "" {
				var jt models.JobType
				if err := tx.Where(models.JobType{Type: job.JobType.Type}).
					FirstOrCreate(&jt).Error; err != nil {
					return fmt.Errorf("job[%d] job_type lookup failed: %w", i, err)
				}
				job.JobTypeID = &jt.ID
				job.JobType = nil
			}

			// ── C. Skills (FirstOrCreate per skill) ───────────────────────────
			var linkedSkills []models.Skill
			for _, s := range job.Skills {
				if s.Skill == "" {
					continue
				}
				var skill models.Skill
				if err := tx.Where(models.Skill{Skill: s.Skill}).
					FirstOrCreate(&skill).Error; err != nil {
					return fmt.Errorf("job[%d] skill '%s' lookup failed: %w", i, s.Skill, err)
				}
				linkedSkills = append(linkedSkills, skill)
			}
			job.Skills = linkedSkills

			// ── D. Geo (lookup only, no create) ───────────────────────────────
			if job.MetaData.GeoData != nil && job.MetaData.GeoData.Province != "" {
				var geo models.GeoData
				err := tx.Where("province = ?", job.MetaData.GeoData.Province).
					First(&geo).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job[%d] province '%s' not registered in geo_data",
						i, job.MetaData.GeoData.Province)
				} else if err != nil {
					return fmt.Errorf("job[%d] geo lookup failed: %w", i, err)
				}
				job.MetaData.GeoDataID = &geo.ID
				job.MetaData.GeoData = nil
			}

			// ── E. Industry (lookup only, no create) ──────────────────────────
			if job.MetaData.Industry != nil && job.MetaData.Industry.Name != "" {
				var industry models.Industry
				err := tx.Where("name = ?", job.MetaData.Industry.Name).
					First(&industry).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job[%d] industry '%s' not registered",
						i, job.MetaData.Industry.Name)
				} else if err != nil {
					return fmt.Errorf("job[%d] industry lookup failed: %w", i, err)
				}
				job.MetaData.IndustryID = &industry.ID
				job.MetaData.Industry = nil
			}

			// ── F. Occupation (lookup only, no create) ────────────────────────
			if job.MetaData.Occupation != nil && job.MetaData.Occupation.Name != "" {
				var occupation models.Occupation
				err := tx.Where("name = ?", job.MetaData.Occupation.Name).
					First(&occupation).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job[%d] occupation '%s' not registered",
						i, job.MetaData.Occupation.Name)
				} else if err != nil {
					return fmt.Errorf("job[%d] occupation lookup failed: %w", i, err)
				}
				job.MetaData.OccupationID = &occupation.ID
				job.MetaData.Occupation = nil
			}

			// ── G. EducationLevel (lookup only, no create) ────────────────────
			if job.MetaData.EducationLevel != nil && job.MetaData.EducationLevel.Level != "" {
				var edu models.EducationLevel
				err := tx.Where("level = ?", job.MetaData.EducationLevel.Level).
					First(&edu).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job[%d] education_level '%s' not registered",
						i, job.MetaData.EducationLevel.Level)
				} else if err != nil {
					return fmt.Errorf("job[%d] education_level lookup failed: %w", i, err)
				}
				job.MetaData.EducationLevelID = &edu.ID
				job.MetaData.EducationLevel = nil
			}

			// ── H. Experience (lookup only, no create) ────────────────────────
			if job.MetaData.Experience != nil && job.MetaData.Experience.Name != "" {
				var exp models.Experience
				err := tx.Where("name = ?", job.MetaData.Experience.Name).
					First(&exp).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("job[%d] experience '%s' not registered",
						i, job.MetaData.Experience.Name)
				} else if err != nil {
					return fmt.Errorf("job[%d] experience lookup failed: %w", i, err)
				}
				job.MetaData.ExperienceID = &exp.ID
				job.MetaData.Experience = nil
			}

			// ── I. Source (FirstOrCreate) ─────────────────────────────────────
			if job.MetaData.Source != nil && job.MetaData.Source.Source != "" {
				var source models.Source
				if err := tx.Where(models.Source{Source: job.MetaData.Source.Source}).
					FirstOrCreate(&source).Error; err != nil {
					return fmt.Errorf("job[%d] source lookup failed: %w", i, err)
				}
				job.MetaData.SourceID = &source.ID
				job.MetaData.Source = nil
			}

			// ── J. AiVersion (FirstOrCreate) ──────────────────────────────────
			if job.MetaData.AiVersion != nil && job.MetaData.AiVersion.Version != "" {
				var av models.AiVersion
				if err := tx.Where(models.AiVersion{Version: job.MetaData.AiVersion.Version}).
					FirstOrCreate(&av).Error; err != nil {
					return fmt.Errorf("job[%d] ai_version lookup failed: %w", i, err)
				}
				job.MetaData.AiVersionID = &av.ID
				job.MetaData.AiVersion = nil
			}

			// ── K. Save JobPost (cascades MetaData + join table) ──────────────
			if err := tx.Create(job).Error; err != nil {
				return fmt.Errorf("job[%d] insert failed: %w", i, err)
			}
			generatedJobIDs[i] = job.ID
		}

		// ── Map true DB IDs onto LSH records ──────────────────────────────────
		for idx := range lshIndexRecords {
			jobGroupIndex := idx / 8
			trueID, exists := generatedJobIDs[jobGroupIndex]
			if !exists {
				return fmt.Errorf("lsh_index[%d] has no matching generated job ID", idx)
			}
			lshIndexRecords[idx].JobPostID = trueID
		}

		// ── Bulk insert LSH records ───────────────────────────────────────────
		if len(lshIndexRecords) > 0 {
			if err := tx.Omit("JobPost").CreateInBatches(&lshIndexRecords, 500).Error; err != nil {
				return fmt.Errorf("lsh_index bulk insert failed: %w", err)
			}
		}

		return nil
	})
}

func (r *JobRepository) BatchUpdateDuplicateJobs(updates []models.JobMetaData) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, up := range updates {
			err := tx.Table("meta_data").
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


//job delete API
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

//get active jobs by industry, province, experience, job type
func (r *JobRepository) GetActiveJobs(industryID, geoDataID, jobTypeID, experienceID *uint) ([]models.JobPost, error) {
	var jobs []models.JobPost

	query := r.withFullPreloads().
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.end_date IS NULL")

	if industryID != nil {
		query = query.Where("meta_data.industry_id = ?", *industryID)
	}
	if geoDataID != nil {
		query = query.Where("meta_data.geo_data_id = ?", *geoDataID)
	}
	if jobTypeID != nil {
		query = query.Where("job_post.job_type_id = ?", *jobTypeID)
	}
	if experienceID != nil {
		query = query.Where("meta_data.experience_id = ?", *experienceID)
	}

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) withFullPreloads() *gorm.DB {
	return r.db.
		Preload("Employer").
		Preload("JobType").
		Preload("Skills").
		Preload("MetaData").
		Preload("MetaData.AiVersion").
		Preload("MetaData.EducationLevel").
		Preload("MetaData.GeoData").
		Preload("MetaData.Industry").
		Preload("MetaData.Occupation").
		Preload("MetaData.Source").
		Preload("MetaData.Experience").
		Preload("MetaData.CrawlerRun")
}

func (r *JobRepository) GetAllIndustries() ([]models.Industry, error) {
	var industries []models.Industry
	if err := r.db.Find(&industries).Error; err != nil {
		return nil, err
	}
	return industries, nil
}

func (r *JobRepository) GetAllExperiences() ([]models.Experience, error) {
	var experiences []models.Experience
	if err := r.db.Find(&experiences).Error; err != nil {
		return nil, err
	}
	return experiences, nil
}

func (r *JobRepository) GetAllProvinces() ([]models.GeoData, error) {
	var provinces []models.GeoData
	if err := r.db.Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *JobRepository) GetAllJobTypes() ([]models.JobType, error) {
	var jobTypes []models.JobType
	if err := r.db.Find(&jobTypes).Error; err != nil {
		return nil, err
	}
	return jobTypes, nil
}

func (r *JobRepository) GetUniqueSkillsCountByIndustry(industryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Skill{}).
		Distinct("skills.id").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.end_date IS NULL", industryID).
		Count(&count).Error
	return count, err
}

func (r *JobRepository) GetMostDemandingSkillByIndustry(industryID uint) (models.SkillDemand, error) {
	var result models.SkillDemand
	err := r.db.Table("skills").
		Select("skills.id, skills.skill, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.end_date IS NULL", industryID).
		Group("skills.id, skills.skill").
		Order("open_job_count DESC").
		Limit(1).
		Scan(&result).Error
	return result, err
}

func (r *JobRepository) GetTop15SkillsByIndustry(industryID uint) ([]models.SkillDemand, error) {
	var results []models.SkillDemand
	err := r.db.Table("skills").
		Select("skills.id, skills.skill, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.end_date IS NULL", industryID).
		Group("skills.id, skills.skill").
		Order("open_job_count DESC").
		Limit(15).
		Scan(&results).Error
	return results, err
}

func (r *JobRepository) GetAllSkillsByIndustry(industryID uint) ([]models.SkillDemand, error) {
	var results []models.SkillDemand
	err := r.db.Table("skills").
		Select("skills.id, skills.skill, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.end_date IS NULL", industryID).
		Group("skills.id, skills.skill").
		Order("open_job_count DESC").
		Scan(&results).Error
	return results, err
}

func (r *JobRepository) GetTopHiringEmployersByIndustry(industryID uint) ([]models.EmployerDemand, error) {
	var results []models.EmployerDemand
	err := r.db.Table("employer").
		Select("employer.id, employer.name, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post ON job_post.employer_id = employer.id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.end_date IS NULL", industryID).
		Group("employer.id, employer.name").
		Order("open_job_count DESC").
		Limit(10).
		Scan(&results).Error
	return results, err
}

func (r *JobRepository) GetLastCrawledJobCount() (int64, error) {
	var count int64

	// Get the latest completed crawler run ID first
	var lastRun models.CrawlerRun
	if err := r.db.Where("status = ?", "COMPLETED").
		Order("finished_at DESC").
		First(&lastRun).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	// Count job posts created under that crawler run
	err := r.db.Model(&models.JobMetaData{}).
		Where("crawler_run_id = ?", lastRun.ID).
		Count(&count).Error

	return count, err
}

func (r *JobRepository) GetTimeSinceLastCrawl() (models.CrawlTimeGap, error) {
	var lastRun models.CrawlerRun

	err := r.db.Where("status = ?", "COMPLETED").
		Order("finished_at DESC").
		First(&lastRun).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CrawlTimeGap{}, nil
	}
	if err != nil {
		return models.CrawlTimeGap{}, err
	}

	now := time.Now()
	gap := now.Sub(*lastRun.FinishedAt)

	return models.CrawlTimeGap{
		LastCrawledAt: lastRun.FinishedAt,
		GapSeconds:    gap.Seconds(),
		GapHuman:      formatDuration(gap),
	}, nil
}

func (r *JobRepository) GetSourcesWithActiveJobCount() ([]models.SourceJobCount, error) {
	var results []models.SourceJobCount

	err := r.db.Table("source").
		Select("source.id, source.source, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN meta_data ON meta_data.source_id = source.id").
		Joins("JOIN job_post ON job_post.id = meta_data.job_post_id").
		Where("meta_data.end_date IS NULL").
		Group("source.id, source.source").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetAllCrawlerRuns() ([]models.CrawlerRun, error) {
	var runs []models.CrawlerRun

	err := r.db.Order("created_at DESC").
		Find(&runs).Error

	return runs, err
}

// formatDuration converts a duration into a human readable string
func formatDuration(d time.Duration) string {
	if d.Hours() >= 24 {
		days := int(d.Hours()) / 24
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%d day(s) %d hour(s) ago", days, hours)
	}
	if d.Hours() >= 1 {
		return fmt.Sprintf("%d hour(s) %d minute(s) ago", int(d.Hours()), int(d.Minutes())%60)
	}
	if d.Minutes() >= 1 {
		return fmt.Sprintf("%d minute(s) %d second(s) ago", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%d second(s) ago", int(d.Seconds()))
}

func (r *JobRepository) GetActiveJobCountWithTrend() (models.JobCountWithTrend, error) {
	var currentCount int64
	var lastMonthCount int64

	now := time.Now()
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayLastMonth := firstDayThisMonth.AddDate(0, -1, 0)

	// Current active job count
	if err := r.db.Model(&models.JobMetaData{}).
		Where("end_date IS NULL").
		Count(&currentCount).Error; err != nil {
		return models.JobCountWithTrend{}, err
	}

	// Last month active job count — jobs that were active during last month window
	if err := r.db.Model(&models.JobMetaData{}).
		Where("posted_at >= ? AND posted_at < ?", firstDayLastMonth, firstDayThisMonth).
		Count(&lastMonthCount).Error; err != nil {
		return models.JobCountWithTrend{}, err
	}

	var changePercent float64
	var trend string

	if lastMonthCount == 0 {
		if currentCount > 0 {
			changePercent = 100.0
			trend = "up"
		} else {
			changePercent = 0.0
			trend = "stable"
		}
	} else {
		changePercent = float64(currentCount-lastMonthCount) / float64(lastMonthCount) * 100
		if changePercent > 0 {
			trend = "up"
		} else if changePercent < 0 {
			trend = "down"
		} else {
			trend = "stable"
		}
	}

	return models.JobCountWithTrend{
		ActiveJobCount: currentCount,
		LastMonthCount: lastMonthCount,
		ChangePercent:  math.Round(changePercent*100) / 100,
		Trend:          trend,
	}, nil
}

func (r *JobRepository) GetActiveJobCountByOccupation() ([]models.OccupationJobCount, error) {
	var results []models.OccupationJobCount

	err := r.db.Table("occupation").
		Select("occupation.id, occupation.name, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.occupation_id = occupation.id AND meta_data.end_date IS NULL").
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("occupation.id, occupation.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetActiveJobCountByIndustry() ([]models.IndustryJobCount, error) {
	var results []models.IndustryJobCount

	err := r.db.Table("industry").
		Select("industry.id, industry.name, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.industry_id = industry.id AND meta_data.end_date IS NULL").
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("industry.id, industry.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetActiveJobCountByExperience() ([]models.ExperienceJobCount, error) {
	var results []models.ExperienceJobCount

	err := r.db.Table("experience").
		Select("experience.id, experience.name, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.experience_id = experience.id AND meta_data.end_date IS NULL").
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("experience.id, experience.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetActiveJobCountByEducationLevel() ([]models.EducationLevelJobCount, error) {
	var results []models.EducationLevelJobCount

	err := r.db.Table("education_level").
		Select("education_level.id, education_level.level, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.education_level_id = education_level.id AND meta_data.end_date IS NULL").
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("education_level.id, education_level.level").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetRemoteVsOnSiteCount() (models.RemoteOnSiteCount, error) {
	var result models.RemoteOnSiteCount

	type row struct {
		IsRemote bool
		Count    int64
	}
	var rows []row

	err := r.db.Table("job_post").
		Select("job_post.is_remote, COUNT(DISTINCT job_post.id) AS count").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.end_date IS NULL").
		Group("job_post.is_remote").
		Scan(&rows).Error

	if err != nil {
		return models.RemoteOnSiteCount{}, err
	}

	for _, r := range rows {
		if r.IsRemote {
			result.RemoteCount = r.Count
		} else {
			result.OnSiteCount = r.Count
		}
	}

	return result, nil
}

func (r *JobRepository) GetActiveJobCountByJobType() ([]models.JobTypeJobCount, error) {
	var results []models.JobTypeJobCount

	err := r.db.Table("job_type").
		Select("job_type.id, job_type.type, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post ON job_post.job_type_id = job_type.id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.end_date IS NULL").
		Group("job_type.id, job_type.type").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetYearlyJobTrendByOccupation(occupationID uint) ([]models.OccupationYearlyTrend, error) {
	var results []models.OccupationYearlyTrend

	now := time.Now()
	startOfTwoYearsAgo := time.Date(now.Year()-2, time.January, 1, 0, 0, 0, 0, now.Location())

	err := r.db.Table("job_post").
		Select("EXTRACT(YEAR FROM meta_data.posted_at)::int AS year, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.occupation_id = ? AND meta_data.posted_at IS NOT NULL AND meta_data.posted_at >= ?", occupationID, startOfTwoYearsAgo).
		Group("EXTRACT(YEAR FROM meta_data.posted_at)").
		Order("year ASC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetTop3JobRolesByOccupation(occupationID uint) ([]models.TopJobRole, error) {
	var results []models.TopJobRole

	err := r.db.Table("job_post").
		Select("job_post.job_role, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.occupation_id = ? AND meta_data.end_date IS NULL", occupationID).
		Group("job_post.job_role").
		Order("open_job_count DESC").
		Limit(3).
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetYearlyJobTrendByIndustry(industryID uint) ([]models.IndustryYearlyTrend, error) {
	var results []models.IndustryYearlyTrend

	now := time.Now()
	startOfTwoYearsAgo := time.Date(now.Year()-2, time.January, 1, 0, 0, 0, 0, now.Location())

	err := r.db.Table("job_post").
		Select("EXTRACT(YEAR FROM meta_data.posted_at)::int AS year, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND meta_data.posted_at IS NOT NULL AND meta_data.posted_at >= ?", industryID, startOfTwoYearsAgo).
		Group("EXTRACT(YEAR FROM meta_data.posted_at)").
		Order("year ASC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetJobCountByExperienceForIndustryAndYear(industryID uint, year int) ([]models.ExperienceYearlyJobCount, error) {
	var results []models.ExperienceYearlyJobCount

	err := r.db.Table("experience").
		Select("experience.id, experience.name, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.experience_id = experience.id AND meta_data.industry_id = ? AND EXTRACT(YEAR FROM meta_data.posted_at) = ?", industryID, year).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("experience.id, experience.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetProvinceWiseJobCountForIndustryAndYear(industryID uint, year int) ([]models.ProvinceJobCount, error) {
	var results []models.ProvinceJobCount

	err := r.db.Table("geo_data").
		Select("geo_data.id, geo_data.province, geo_data.latitude, geo_data.longitude, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.geo_data_id = geo_data.id AND meta_data.industry_id = ? AND EXTRACT(YEAR FROM meta_data.posted_at) = ?", industryID, year).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("geo_data.id, geo_data.province, geo_data.latitude, geo_data.longitude").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetJobCountByEducationLevelForIndustryAndYear(industryID uint, year int) ([]models.EducationLevelYearlyJobCount, error) {
	var results []models.EducationLevelYearlyJobCount

	err := r.db.Table("education_level").
		Select("education_level.id, education_level.level, COALESCE(COUNT(DISTINCT job_post.id), 0) AS open_job_count").
		Joins("LEFT JOIN meta_data ON meta_data.education_level_id = education_level.id AND meta_data.industry_id = ? AND EXTRACT(YEAR FROM meta_data.posted_at) = ?", industryID, year).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("education_level.id, education_level.level").
		Order("open_job_count DESC").
		Scan(&results).Error

	return results, err
}

func (r *JobRepository) GetTopHiringEmployersForIndustryAndYear(industryID uint, year int) ([]models.TopEmployerByIndustryYear, error) {
	var results []models.TopEmployerByIndustryYear

	err := r.db.Table("employer").
		Select("employer.id, employer.name, COUNT(DISTINCT job_post.id) AS open_job_count").
		Joins("JOIN job_post ON job_post.employer_id = employer.id").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.industry_id = ? AND EXTRACT(YEAR FROM meta_data.posted_at) = ?", industryID, year).
		Group("employer.id, employer.name").
		Order("open_job_count DESC").
		Limit(10).
		Scan(&results).Error

	return results, err
}