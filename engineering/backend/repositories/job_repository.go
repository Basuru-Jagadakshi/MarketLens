package repositories

import (
	"errors"
	"marketlens-go-backend/models"

	"gorm.io/gorm"
)



type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(input models.JobPost) (models.JobPost, error) {
	
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

		return nil
	})

	if err != nil {
		return models.JobPost{}, err
	}

	var completeJob models.JobPost
	r.db.Preload("JobType").Preload("MetaData.Geo").Preload("Skills").First(&completeJob, input.ID)

	return completeJob, nil
}


func (r *JobRepository) GetAllJobs() ([]models.JobPost, error) {
	
	var jobs []models.JobPost

	err := r.db.Preload("JobType").Preload("MetaData.Geo").Preload("Skills").Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}



func (r *JobRepository) UpdateJob(id uint, input models.JobPost) (models.JobPost, error) {
	
	var existingJob models.JobPost
	
	if err := r.db.Preload("MetaData").First(&existingJob, id).Error; err != nil {
		return models.JobPost{}, err
	}

	input.ID = id
	input.MetaData.JobPostID = id
	input.MetaData.ID = existingJob.MetaData.ID

	err := r.db.Transaction(func(tx *gorm.DB) error {
		
		if input.MetaData.Geo != nil && input.MetaData.Geo.Province != "" {
			var existingGeo models.GeoData
			err := tx.Where("province = ?", input.MetaData.Geo.Province).First(&existingGeo).Error
			if err != nil {
				return errors.New("provided province does not match system geo registers")
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

		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&input).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.JobPost{}, err
	}

	var completeJob models.JobPost
	r.db.Preload("JobType").Preload("MetaData.Geo").Preload("Skills").First(&completeJob, id)

	return completeJob, nil
}