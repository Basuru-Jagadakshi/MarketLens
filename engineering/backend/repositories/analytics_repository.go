package repositories

import (
	"fmt"
	"marketlens-go-backend/models"
	"time"

	"gorm.io/gorm"
)

// This function builds a subquery of job_post.ids that fall under the given
// occupation/industry hierarchy level and id — used to scope other aggregations
// (like employment sector breakdown) to that level.
func (r *JobRepository) buildJobPostIDsForLevel(standard, level string, id uint, fromDate, toDate time.Time) (*gorm.DB, error) {
	toDateExclusive := toDate.AddDate(0, 0, 1)

	if standard == "occupation" {
		switch level {
		case "major-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Joins("JOIN minor_group ON minor_group.id = unit_group.minor_group_id").
				Joins("JOIN sub_major_group ON sub_major_group.id = minor_group.sub_major_group_id").
				Where("sub_major_group.major_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "sub-major-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Joins("JOIN minor_group ON minor_group.id = unit_group.minor_group_id").
				Where("minor_group.sub_major_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "minor-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Where("unit_group.minor_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "unit-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Where("occupation_group.unit_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "occupation-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Where("meta_data.occupation_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		default:
			return nil, fmt.Errorf("invalid level '%s' for standard 'occupation'", level)
		}
	} else if standard == "industry" {
		switch level {
		case "industry-sector":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Joins("JOIN industry_group ON industry_group.id = industry_class.industry_group_id").
				Joins("JOIN industry_division ON industry_division.id = industry_group.industry_division_id").
				Where("industry_division.industry_sector_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "industry-division":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Joins("JOIN industry_group ON industry_group.id = industry_class.industry_group_id").
				Where("industry_group.industry_division_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "industry-group":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Where("industry_class.industry_group_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "industry-class":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Where("industry_subclass.industry_class_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		case "industry-subclass":
			return r.db.Table("job_post").
				Select("job_post.id").
				Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
				Where("meta_data.industry_subclass_id = ?", id).
				Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive), nil

		default:
			return nil, fmt.Errorf("invalid level '%s' for standard 'industry'", level)
		}
	}

	return nil, fmt.Errorf("invalid standard '%s', must be 'occupation' or 'industry'", standard)
}

// This function returns the top hiring employers (by summed vacancy count) for jobs
// under the given occupation hierarchy level and id, filtered by date range.
func (r *JobRepository) GetTopHiringEmployersByOccupationLevel(level string, id uint, fromDate, toDate time.Time) ([]models.EmployerDemand, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel("occupation", level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var results []models.EmployerDemand
	err = r.db.Table("employer").
		Select("employer.id, employer.name, SUM(job_post.no_of_vacancies) AS open_job_count").
		Joins("JOIN job_post ON job_post.employer_id = employer.id").
		Where("job_post.id IN (?)", jobPostIDs).
		Group("employer.id, employer.name").
		Order("open_job_count DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query top hiring employers for occupation/%s/%d: %w", level, id, err)
	}

	return results, nil
}

// This function returns all skills (paginated, by summed vacancy count) for jobs
// under the given occupation hierarchy level and id, filtered by date range.
func (r *JobRepository) GetAllSkillsByOccupationLevel(level string, id uint, fromDate, toDate time.Time, limit, offset int) ([]models.SkillDemand, int64, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel("occupation", level, id, fromDate, toDate)
	if err != nil {
		return nil, 0, err
	}

	baseQuery := r.db.Table("skills").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Where("job_post.id IN (?)", jobPostIDs)

	var total int64
	if err := baseQuery.Session(&gorm.Session{}).
		Distinct("skills.id").
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count skills for occupation/%s/%d: %w", level, id, err)
	}

	var results []models.SkillDemand
	query := baseQuery.Session(&gorm.Session{}).
		Select("skills.id, skills.skill, SUM(job_post.no_of_vacancies) AS open_job_count").
		Group("skills.id, skills.skill").
		Order("open_job_count DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query all skills for occupation/%s/%d: %w", level, id, err)
	}

	return results, total, nil
}

// This function returns the top 15 skills (by summed vacancy count) for jobs
// under the given occupation hierarchy level and id, filtered by date range.
func (r *JobRepository) GetTop15SkillsByOccupationLevel(level string, id uint, fromDate, toDate time.Time) ([]models.SkillDemand, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel("occupation", level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var results []models.SkillDemand
	err = r.db.Table("skills").
		Select("skills.id, skills.skill, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Joins("JOIN job_post_skills ON job_post_skills.skill_id = skills.id").
		Joins("JOIN job_post ON job_post.id = job_post_skills.job_post_id").
		Where("job_post.id IN (?)", jobPostIDs).
		Group("skills.id, skills.skill").
		Order("open_job_count DESC").
		Limit(15).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query top 15 skills for occupation/%s/%d: %w", level, id, err)
	}

	return results, nil
}

// This function returns job type breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetJobTypeByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.JobTypeJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var results []models.JobTypeJobCount
	err = r.db.Table("job_type").
		Select("job_type.id, job_type.type, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Joins(
			"LEFT JOIN job_post ON job_post.job_type_id = job_type.id AND job_post.id IN (?)",
			jobPostIDs,
		).
		Group("job_type.id, job_type.type").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query job type breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns remote vs on-site job counts for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetRemoteOnSiteHybridByLevel(standard, level string, id uint, fromDate, toDate time.Time) (models.RemoteOnSiteHybridCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return models.RemoteOnSiteHybridCount{}, err
	}

	type workModeRow struct {
		WorkMode models.WorkMode
		Count    int64
	}
	var rows []workModeRow

	err = r.db.Table("job_post").
		Select("job_post.work_mode, COALESCE(SUM(job_post.no_of_vacancies), 0) AS count").
		Where("job_post.id IN (?)", jobPostIDs).
		Group("job_post.work_mode").
		Scan(&rows).Error

	if err != nil {
		return models.RemoteOnSiteHybridCount{}, fmt.Errorf("failed to query remote/on-site breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	var result models.RemoteOnSiteHybridCount
	for _, row := range rows {
		switch row.WorkMode {
		case models.WorkModeRemote:
			result.RemoteCount = row.Count
		case models.WorkModeOnsite:
			result.OnSiteCount = row.Count
		case models.WorkModeHybrid:
			result.HybridCount = row.Count
		}
	}

	return result, nil
}

// This function returns vocational education breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetVocationalEducationByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.VocationalEducationJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.VocationalEducationJobCount
	err = r.db.Table("vocational_education").
		Select("vocational_education.id, vocational_education.level, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("vocational_education.created_at <= ?", endOfToDate).
		Where("vocational_education.deleted_at IS NULL OR vocational_education.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.vocational_education_id = vocational_education.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("vocational_education.id, vocational_education.level").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query vocational education breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns gender breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetGenderByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.GenderJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.GenderJobCount
	err = r.db.Table("gender").
		Select("gender.id, gender.gender_type, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("gender.created_at <= ?", endOfToDate).
		Where("gender.deleted_at IS NULL OR gender.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.gender_id = gender.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("gender.id, gender.gender_type").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query gender breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns formality breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetFormalityByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.FormalityJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.FormalityJobCount
	err = r.db.Table("formality").
		Select("formality.id, formality.formality_type, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("formality.created_at <= ?", endOfToDate).
		Where("formality.deleted_at IS NULL OR formality.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.formality_id = formality.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("formality.id, formality.formality_type").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query formality breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns education level breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetEducationLevelByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.EducationLevelJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.EducationLevelJobCount
	err = r.db.Table("education_level").
		Select("education_level.id, education_level.level, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("education_level.created_at <= ?", endOfToDate).
		Where("education_level.deleted_at IS NULL OR education_level.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.education_level_id = education_level.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("education_level.id, education_level.level").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query education level breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns province breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetProvinceByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.ProvinceJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var results []models.ProvinceJobCount
	err = r.db.Table("geo_data").
		Select("geo_data.id, geo_data.province, geo_data.latitude, geo_data.longitude, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Joins(
			"LEFT JOIN meta_data ON meta_data.geo_data_id = geo_data.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("geo_data.id, geo_data.province, geo_data.latitude, geo_data.longitude").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query province breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns experience-level breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetExperienceByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.ExperienceJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.ExperienceJobCount
	err = r.db.Table("experience").
		Select("experience.id, experience.name, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("experience.created_at <= ?", endOfToDate).
		Where("experience.deleted_at IS NULL OR experience.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.experience_id = experience.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("experience.id, experience.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query experience breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns employment sector breakdown (with job counts) for the given
// occupation/industry hierarchy level and id, filtered by date range.
func (r *JobRepository) GetEmploymentSectorByLevel(standard, level string, id uint, fromDate, toDate time.Time) ([]models.EmploymentSectorJobCount, error) {
	jobPostIDs, err := r.buildJobPostIDsForLevel(standard, level, id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var results []models.EmploymentSectorJobCount
	err = r.db.Table("employment_sector").
		Select("employment_sector.id, employment_sector.sector, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("employment_sector.created_at <= ?", endOfToDate).
		Where("employment_sector.deleted_at IS NULL OR employment_sector.deleted_at >= ?", fromDate).
		Joins(
			"LEFT JOIN meta_data ON meta_data.employment_sector_id = employment_sector.id "+
				"AND meta_data.job_post_id IN (?)",
			jobPostIDs,
		).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("employment_sector.id, employment_sector.sector").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query employment sector breakdown for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, nil
}

// This function returns the immediate child-level entities under a given occupation/industry
// hierarchy level and id, each with its aggregated job count for the given date range.
func (r *JobRepository) GetLevelChildren(standard, level string, id uint, fromDate, toDate time.Time) ([]models.LevelChildJobCount, string, error) {
	var results []models.LevelChildJobCount
	var childLevel string
	var query *gorm.DB

	toDateExclusive := toDate.AddDate(0, 0, 1)

	if standard == "occupation" {
		switch level {
		case "major-group":
			childLevel = "sub-major-group"
			query = r.db.Table("sub_major_group").
				Select("sub_major_group.id, sub_major_group.name, sub_major_group.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("sub_major_group.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("sub_major_group.deleted_at IS NULL OR sub_major_group.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN minor_group ON minor_group.sub_major_group_id = sub_major_group.id").
				Joins("LEFT JOIN unit_group ON unit_group.minor_group_id = minor_group.id").
				Joins("LEFT JOIN occupation_group ON occupation_group.unit_group_id = unit_group.id").
				Joins("LEFT JOIN meta_data ON meta_data.occupation_group_id = occupation_group.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("sub_major_group.major_group_id = ?", id).
				Group("sub_major_group.id, sub_major_group.name, sub_major_group.code")

		case "sub-major-group":
			childLevel = "minor-group"
			query = r.db.Table("minor_group").
				Select("minor_group.id, minor_group.name, minor_group.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("minor_group.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("minor_group.deleted_at IS NULL OR minor_group.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN unit_group ON unit_group.minor_group_id = minor_group.id").
				Joins("LEFT JOIN occupation_group ON occupation_group.unit_group_id = unit_group.id").
				Joins("LEFT JOIN meta_data ON meta_data.occupation_group_id = occupation_group.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("minor_group.sub_major_group_id = ?", id).
				Group("minor_group.id, minor_group.name, minor_group.code")

		case "minor-group":
			childLevel = "unit-group"
			query = r.db.Table("unit_group").
				Select("unit_group.id, unit_group.name, unit_group.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("unit_group.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("unit_group.deleted_at IS NULL OR unit_group.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN occupation_group ON occupation_group.unit_group_id = unit_group.id").
				Joins("LEFT JOIN meta_data ON meta_data.occupation_group_id = occupation_group.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("unit_group.minor_group_id = ?", id).
				Group("unit_group.id, unit_group.name, unit_group.code")

		case "unit-group":
			childLevel = "occupation-group"
			query = r.db.Table("occupation_group").
				Select("occupation_group.id, occupation_group.name, occupation_group.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("occupation_group.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("occupation_group.deleted_at IS NULL OR occupation_group.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN meta_data ON meta_data.occupation_group_id = occupation_group.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("occupation_group.unit_group_id = ?", id).
				Group("occupation_group.id, occupation_group.name, occupation_group.code")

		case "occupation-group":
			return nil, "", fmt.Errorf("'occupation-group' is a leaf level and has no children")

		default:
			return nil, "", fmt.Errorf("invalid level '%s' for standard 'occupation'", level)
		}
	} else if standard == "industry" {
		switch level {
		case "industry-sector":
			childLevel = "industry-division"
			query = r.db.Table("industry_division").
				Select("industry_division.id, industry_division.name, industry_division.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("industry_division.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("industry_division.deleted_at IS NULL OR industry_division.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN industry_group ON industry_group.industry_division_id = industry_division.id").
				Joins("LEFT JOIN industry_class ON industry_class.industry_group_id = industry_group.id").
				Joins("LEFT JOIN industry_subclass ON industry_subclass.industry_class_id = industry_class.id").
				Joins("LEFT JOIN meta_data ON meta_data.industry_subclass_id = industry_subclass.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("industry_division.industry_sector_id = ?", id).
				Group("industry_division.id, industry_division.name, industry_division.code")

		case "industry-division":
			childLevel = "industry-group"
			query = r.db.Table("industry_group").
				Select("industry_group.id, industry_group.name, industry_group.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("industry_group.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("industry_group.deleted_at IS NULL OR industry_group.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN industry_class ON industry_class.industry_group_id = industry_group.id").
				Joins("LEFT JOIN industry_subclass ON industry_subclass.industry_class_id = industry_class.id").
				Joins("LEFT JOIN meta_data ON meta_data.industry_subclass_id = industry_subclass.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("industry_group.industry_division_id = ?", id).
				Group("industry_group.id, industry_group.name, industry_group.code")

		case "industry-group":
			childLevel = "industry-class"
			query = r.db.Table("industry_class").
				Select("industry_class.id, industry_class.name, industry_class.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("industry_class.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("industry_class.deleted_at IS NULL OR industry_class.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN industry_subclass ON industry_subclass.industry_class_id = industry_class.id").
				Joins("LEFT JOIN meta_data ON meta_data.industry_subclass_id = industry_subclass.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("industry_class.industry_group_id = ?", id).
				Group("industry_class.id, industry_class.name, industry_class.code")

		case "industry-class":
			childLevel = "industry-subclass"
			query = r.db.Table("industry_subclass").
				Select("industry_subclass.id, industry_subclass.name, industry_subclass.code, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
				Where("industry_subclass.created_at <= ?", toDate.Add(23*time.Hour+59*time.Minute+59*time.Second)).
				Where("industry_subclass.deleted_at IS NULL OR industry_subclass.deleted_at >= ?", fromDate).
				Joins("LEFT JOIN meta_data ON meta_data.industry_subclass_id = industry_subclass.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
				Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
				Where("industry_subclass.industry_class_id = ?", id).
				Group("industry_subclass.id, industry_subclass.name, industry_subclass.code")

		case "industry-subclass":
			return nil, "", fmt.Errorf("'industry-subclass' is a leaf level and has no children")

		default:
			return nil, "", fmt.Errorf("invalid level '%s' for standard 'industry'", level)
		}
	} else {
		return nil, "", fmt.Errorf("invalid standard '%s', must be 'occupation' or 'industry'", standard)
	}

	err := query.Order("open_job_count DESC").Scan(&results).Error
	if err != nil {
		return nil, "", fmt.Errorf("failed to query children for %s/%s/%d: %w", standard, level, id, err)
	}

	return results, childLevel, nil
}

// This function returns the total job count for a given occupation/industry hierarchy level and id,
// filtered by date range. standard is "occupation" or "industry"; level depends on the standard.
func (r *JobRepository) GetTotalJobCountByLevel(standard, level string, id uint, fromDate, toDate time.Time) (int64, error) {
	toDateExclusive := toDate.AddDate(0, 0, 1)

	query := r.db.Table("job_post").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive)

	if standard == "occupation" {
		switch level {
		case "major-group":
			query = query.
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Joins("JOIN minor_group ON minor_group.id = unit_group.minor_group_id").
				Joins("JOIN sub_major_group ON sub_major_group.id = minor_group.sub_major_group_id").
				Where("sub_major_group.major_group_id = ?", id)

		case "sub-major-group":
			query = query.
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Joins("JOIN minor_group ON minor_group.id = unit_group.minor_group_id").
				Where("minor_group.sub_major_group_id = ?", id)

		case "minor-group":
			query = query.
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Joins("JOIN unit_group ON unit_group.id = occupation_group.unit_group_id").
				Where("unit_group.minor_group_id = ?", id)

		case "unit-group":
			query = query.
				Joins("JOIN occupation_group ON occupation_group.id = meta_data.occupation_group_id").
				Where("occupation_group.unit_group_id = ?", id)

		case "occupation-group":
			query = query.Where("meta_data.occupation_group_id = ?", id)

		default:
			return 0, fmt.Errorf("invalid level '%s' for standard 'occupation'", level)
		}
	} else if standard == "industry" {
		switch level {
		case "industry-sector":
			query = query.
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Joins("JOIN industry_group ON industry_group.id = industry_class.industry_group_id").
				Joins("JOIN industry_division ON industry_division.id = industry_group.industry_division_id").
				Where("industry_division.industry_sector_id = ?", id)

		case "industry-division":
			query = query.
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Joins("JOIN industry_group ON industry_group.id = industry_class.industry_group_id").
				Where("industry_group.industry_division_id = ?", id)

		case "industry-group":
			query = query.
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Joins("JOIN industry_class ON industry_class.id = industry_subclass.industry_class_id").
				Where("industry_class.industry_group_id = ?", id)

		case "industry-class":
			query = query.
				Joins("JOIN industry_subclass ON industry_subclass.id = meta_data.industry_subclass_id").
				Where("industry_subclass.industry_class_id = ?", id)

		case "industry-subclass":
			query = query.Where("meta_data.industry_subclass_id = ?", id)

		default:
			return 0, fmt.Errorf("invalid level '%s' for standard 'industry'", level)
		}
	} else {
		return 0, fmt.Errorf("invalid standard '%s', must be 'occupation' or 'industry'", standard)
	}

	var total int64
	err := query.Select("COALESCE(SUM(job_post.no_of_vacancies), 0)").Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("failed to query total job count for %s/%s/%d: %w", standard, level, id, err)
	}

	return total, nil
}

// This function returns vacancy counts grouped by industry sector, for jobs posted within the given date range
func (r *JobRepository) GetIndustryJobCountByDateRange(fromDate, toDate time.Time) ([]models.IndustryJobCount, error) {
	var results []models.IndustryJobCount
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	toDateExclusive := toDate.AddDate(0, 0, 1)

	err := r.db.Table("industry_sector").
		Select("industry_sector.id, industry_sector.name, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("industry_sector.created_at <= ?", endOfToDate).
		Where("industry_sector.deleted_at IS NULL OR industry_sector.deleted_at >= ?", fromDate).
		Joins("LEFT JOIN industry_division ON industry_division.industry_sector_id = industry_sector.id").
		Joins("LEFT JOIN industry_group ON industry_group.industry_division_id = industry_division.id").
		Joins("LEFT JOIN industry_class ON industry_class.industry_group_id = industry_group.id").
		Joins("LEFT JOIN industry_subclass ON industry_subclass.industry_class_id = industry_class.id").
		Joins("LEFT JOIN meta_data ON meta_data.industry_subclass_id = industry_subclass.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("industry_sector.id, industry_sector.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query industry job count by date range: %w", err)
	}

	return results, nil
}

// This function returns vacancy counts grouped by major group (occupation), for jobs posted within the given date range
func (r *JobRepository) GetOccupationJobCountByDateRange(fromDate, toDate time.Time) ([]models.OccupationJobCount, error) {
	var results []models.OccupationJobCount
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	toDateExclusive := toDate.AddDate(0, 0, 1)

	err := r.db.Table("major_group").
		Select("major_group.id, major_group.name, COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count").
		Where("major_group.created_at <= ?", endOfToDate).
		Where("major_group.deleted_at IS NULL OR major_group.deleted_at >= ?", fromDate).
		Joins("LEFT JOIN sub_major_group ON sub_major_group.major_group_id = major_group.id").
		Joins("LEFT JOIN minor_group ON minor_group.sub_major_group_id = sub_major_group.id").
		Joins("LEFT JOIN unit_group ON unit_group.minor_group_id = minor_group.id").
		Joins("LEFT JOIN occupation_group ON occupation_group.unit_group_id = unit_group.id").
		Joins("LEFT JOIN meta_data ON meta_data.occupation_group_id = occupation_group.id AND meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
		Joins("LEFT JOIN job_post ON job_post.id = meta_data.job_post_id").
		Group("major_group.id, major_group.name").
		Order("open_job_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query occupation job count by date range: %w", err)
	}

	return results, nil
}

// This function returns the total sum of no_of_vacancies for jobs posted within the given date range
func (r *JobRepository) GetTotalVacancyCount(fromDate, toDate time.Time) (int64, error) {
	var total int64
	toDateExclusive := toDate.AddDate(0, 0, 1)

	err := r.db.Table("job_post").
		Select("COALESCE(SUM(job_post.no_of_vacancies), 0)").
		Joins("JOIN meta_data ON meta_data.job_post_id = job_post.id").
		Where("meta_data.posted_at >= ? AND meta_data.posted_at < ?", fromDate, toDateExclusive).
		Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("failed to query total vacancy count: %w", err)
	}

	return total, nil
}

// Get vacancy trend for a selected date range
// This function returns job counts grouped day by day within the given date range
func (r *JobRepository) GetVacancyTrendDaily(fromDate, toDate time.Time) ([]models.VacancyTrendPoint, error) {
	var rows []struct {
		PeriodStart  time.Time
		OpenJobCount int64
	}

	err := r.db.Raw(`
		SELECT day_series.day AS period_start,
		       COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count
		FROM generate_series(?::date, ?::date, '1 day'::interval) AS day_series(day)
		LEFT JOIN meta_data ON meta_data.posted_at >= day_series.day
		                    AND meta_data.posted_at < day_series.day + interval '1 day'
		LEFT JOIN job_post ON job_post.id = meta_data.job_post_id
		GROUP BY day_series.day
		ORDER BY day_series.day ASC
	`, fromDate, toDate).Scan(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query daily vacancy trend: %w", err)
	}

	results := make([]models.VacancyTrendPoint, 0, len(rows))
	for _, row := range rows {
		results = append(results, models.VacancyTrendPoint{
			Label:        formatWeekLabel(row.PeriodStart),
			OpenJobCount: row.OpenJobCount,
		})
	}

	return results, nil
}

// This function returns job counts grouped month by month within the given date range
func (r *JobRepository) GetVacancyTrendMonthly(fromDate, toDate time.Time) ([]models.VacancyTrendPoint, error) {
	var rows []struct {
		PeriodStart  time.Time
		OpenJobCount int64
	}

	err := r.db.Raw(`
		SELECT month_series.month AS period_start,
		       COALESCE(SUM(job_post.no_of_vacancies), 0) AS open_job_count
		FROM generate_series(
		         date_trunc('month', ?::date),
		         date_trunc('month', ?::date),
		         '1 month'::interval
		     ) AS month_series(month)
		LEFT JOIN meta_data ON meta_data.posted_at >= month_series.month
		                    AND meta_data.posted_at < month_series.month + interval '1 month'
		                    AND meta_data.posted_at >= ?
		                    AND meta_data.posted_at <= ?
		LEFT JOIN job_post ON job_post.id = meta_data.job_post_id
		GROUP BY month_series.month
		ORDER BY month_series.month ASC
	`, fromDate, toDate, fromDate, toDate).Scan(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query monthly vacancy trend: %w", err)
	}

	results := make([]models.VacancyTrendPoint, 0, len(rows))
	for _, row := range rows {
		results = append(results, models.VacancyTrendPoint{
			Label:        row.PeriodStart.Format("Jan 2006"),
			OpenJobCount: row.OpenJobCount,
		})
	}

	n := len(rows)
	if n == 0 {
		return results, nil
	}

	firstBucketStart := rows[0].PeriodStart
	firstBucketNaturalEnd := firstBucketStart.AddDate(0, 1, 0).AddDate(0, 0, -1)
	lastBucketStart := rows[n-1].PeriodStart
	lastBucketNaturalEnd := lastBucketStart.AddDate(0, 1, 0).AddDate(0, 0, -1)

	firstIsPartial := fromDate.After(firstBucketStart)
	lastIsPartial := toDate.Before(lastBucketNaturalEnd)

	switch {
	case n == 1 && firstIsPartial && lastIsPartial:
		results[0].Label = fmt.Sprintf("%s - %s", fromDate.Format("Jan 2"), toDate.Format("Jan 2, 2006"))

	default:
		if firstIsPartial {
			results[0].Label = fmt.Sprintf("%s - %s", fromDate.Format("Jan 2"), firstBucketNaturalEnd.Format("Jan 2, 2006"))
		}
		if lastIsPartial {
			results[n-1].Label = fmt.Sprintf("%s - %s", lastBucketStart.Format("Jan 2"), toDate.Format("Jan 2, 2006"))
		}
	}

	return results, nil
}

func formatWeekLabel(t time.Time) string {
	return t.Format("Jan 2")
}
