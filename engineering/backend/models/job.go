package models

import (
	"time"
	"github.com/lib/pq"
)


type JobType struct {
	ID 		uint 		`json:"id" gorm:"primaryKey"`
	Name 	string 		`json:"name" gorm:"unique;not null"`
}

type GeoData struct {
	ID 			uint 		`json:"id" gorm:"primaryKey"`
	Latitude 	float64 	`json:"lat" gorm:"type:decimal(9,6);not null"`
	Longitude 	float64 	`json:"lng" gorm:"type:decimal(9,6);not null"`
	Province 	string 		`json:"province" gorm:"size:50;not null"`
}

type Skill struct {
	ID 		uint 		`json:"id" gorm:"primaryKey"`
	Name 	string 		`json:"name" gorm:"unique;not null"`
}

type CrawlerRun struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	StartedAt  time.Time  `json:"started_at" gorm:"default:CURRENT_TIMESTAMP"`
	FinishedAt *time.Time `json:"finished_at"`
	Status     string     `json:"status" gorm:"size:20;default:'RUNNING'"` // 'RUNNING', 'COMPLETED', 'FAILED'
}

func (CrawlerRun) TableName() string {
	return "crawler_runs"
}

type JobPost struct {
	ID 								uint 			`json:"id" gorm:"primaryKey"`
	Employer 						string 			`json:"employer" gorm:"size:255;not null"`
	JobRole 						string 			`json:"job_role" gorm:"size:255;not null"`
	JobTypeID 						uint 			`json:"job_type_id"`
	JobType 						JobType 		`json:"job_type" gorm:"foreignKey:JobTypeID"`
	KeyResponsibilities 			string 			`json:"key_responsibilities" gorm:"type:text"`
	Qualifications 					string 			`json:"qualifications" gorm:"type:text"`
	LocationDescription 			string 			`json:"location" gorm:"size:255"`
	Offers 							string 			`json:"offers" gorm:"type:text"`
	IsRemote 						bool 			`json:"is_remote" gorm:"default:false"`
	CreatedAt 						time.Time 		`json:"created_at"`
	MetaData					 	JobMetaData 	`json:"meta_data" gorm:"foreignKey:JobPostID;constraint:OnDelete:CASCADE;"`
	Skills 							[]Skill 		`json:"skills" gorm:"many2many:job_skills;constraint:OnDelete:CASCADE;"`
}

type JobMetaData struct {
	ID 								uint			`json:"-" gorm:"primaryKey"`
	JobPostID 						uint			`json:"job_post_id" gorm:"unique;not null"`
	GeoID 							*uint			`json:"-"`
	Geo 							*GeoData		`json:"geo" gorm:"foreignKey:GeoID"`
	PostedAt 						time.Time		`json:"posted_at" gorm:"default:CURRENT_TIMESTAMP"`
	Source 							string			`json:"source" gorm:"size:100"`
	StandarizedCategory 			string			`json:"standardized_category" gorm:"size:100"`
	Seniority 						string			`json:"seniority" gorm:"size:50"`
	ConfidenceScore 				float64			`json:"confidence_score" gorm:"type:decimal(3,2)"`
	AiVersion 						string			`json:"ai_version" gorm:"size:50"`
	HasError 						bool			`json:"error" gorm:"default:false"`
	EndDate          				*time.Time    	`json:"end_date"`                            
	CrawlerRunID     				*uint         	`json:"crawler_run_id"`                      
	CrawlerRun       				*CrawlerRun   	`json:"-" gorm:"foreignKey:CrawlerRunID"`   
	MinhashSignature 				pq.Int64Array 	`json:"minhash_signature" gorm:"type:integer[]"`
}

func (JobMetaData) TableName() string {
	return "job_metadata"
}

type LshIndex struct {
	BucketKey string   `gorm:"primaryKey;size:64;not null" json:"bucket_key"` 
	BandNo    int      `gorm:"not null" json:"band_no"`                     
	JobPostID uint     `gorm:"primaryKey;not null" json:"job_post_id"`
	JobPost   *JobPost `gorm:"foreignKey:JobPostID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (LshIndex) TableName() string {
	return "lsh_index"
}

type BatchSavePayload struct {
	NewJobs    []JobPost  `json:"new_jobs" binding:"required"`
	LshIndexes []LshIndex `json:"lsh_indexes" binding:"required"`
}

type BatchUpdatePayload struct {
	Duplicates []JobMetaData `json:"duplicates" binding:"required"`
}

type ReconciliationPayload struct {
	CrawlerRunID uint `json:"crawler_run_id" binding:"required"`
}