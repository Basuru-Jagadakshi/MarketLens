package models

import "time"


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
	JobPostID 						uint			`json:"-" gorm:"unique;not null"`
	GeoID 							*uint			`json:"-"`
	Geo 							*GeoData		`json:"geo" gorm:"foreignKey:GeoID"`
	PostedAt 						time.Time		`json:"posted_at"`
	Source 							string			`json:"source" gorm:"size:100"`
	StandarizedCategory 			string			`json:"standardized_category" gorm:"size:100"`
	Seniority 						string			`json:"seniority" gorm:"size:50"`
	ConfidenceScore 				float64			`json:"confidence_score" gorm:"type:decimal(3,2)"`
	AiVersion 						string			`json:"ai_version" gorm:"size:50"`
	HasError 						bool			`json:"error" gorm:"default:false"`
}



func (JobMetaData) TableName() string {
	return "job_metadata"
}