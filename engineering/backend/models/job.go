package models

import (
	"time"
	"github.com/lib/pq"
)

// ── Lookup / reference tables ─────────────────────────────────────────────────

type Employer struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Name      string    `json:"name"       gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Employer) TableName() string { return "employer" }

type JobType struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Type      string    `json:"type"       gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (JobType) TableName() string       { return "job_type" }

type Skill struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Skill     string    `json:"skill"      gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Skill) TableName() string         { return "skills" }

type AiVersion struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Version   string    `json:"version"    gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AiVersion) TableName() string { return "ai_version" }

type EducationLevel struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Level     string    `json:"level"      gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EducationLevel) TableName() string { return "education_level" }

type Industry struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Name      string    `json:"name"       gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Industry) TableName() string { return "industry" }

type Occupation struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Name      string    `json:"name"       gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Occupation) TableName() string { return "occupation" }

type Source struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Source    string    `json:"source"     gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Source) TableName() string { return "source" }

type Experience struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Name      string    `json:"name"       gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Experience) TableName() string { return "experience" }

type GeoData struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Longitude float64   `json:"lng"        gorm:"type:decimal(9,6)"`
	Latitude  float64   `json:"lat"        gorm:"type:decimal(9,6)"`
	Province  string    `json:"province"   gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GeoData) TableName() string { return "geo_data" }

// ── Crawler run ───────────────────────────────────────────────────────────────

type CrawlerRun struct {
	ID         uint       `json:"id"          gorm:"primaryKey"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Status     string     `json:"status"      gorm:"size:50"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (CrawlerRun) TableName() string { return "crawler_runs" }

// ── Core job tables ───────────────────────────────────────────────────────────

type JobPost struct {
	ID             uint       `json:"id"              gorm:"primaryKey"`
	EmployerID     *uint      `json:"employer_id"`
	Employer       *Employer  `json:"employer"        gorm:"foreignKey:EmployerID"`
	JobTypeID      *uint      `json:"job_type_id"`
	JobType        *JobType   `json:"job_type"        gorm:"foreignKey:JobTypeID"`
	JobRole        string     `json:"job_role"        gorm:"size:255;not null"`
	IsRemote       bool       `json:"is_remote"       gorm:"default:false"`
	JobDescription string     `json:"job_description" gorm:"type:text"`
	Location       string     `json:"location"        gorm:"size:255"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Associations
	MetaData       JobMetaData `json:"meta_data" gorm:"foreignKey:JobPostID;constraint:OnDelete:CASCADE;"`
	Skills         []Skill     `json:"skills"    gorm:"many2many:job_post_skills;constraint:OnDelete:CASCADE;"`
}

func (JobPost) TableName() string { return "job_post" }

type JobMetaData struct {
	ID               uint            `json:"-"                  gorm:"primaryKey"`
	JobPostID        uint            `json:"job_post_id"        gorm:"unique;not null"`
	AiVersionID      *uint           `json:"ai_version_id"`
	AiVersion        *AiVersion      `json:"ai_version"         gorm:"foreignKey:AiVersionID"`
	EducationLevelID *uint           `json:"education_level_id"`
	EducationLevel   *EducationLevel `json:"education_level"    gorm:"foreignKey:EducationLevelID"`
	CrawlerRunID     *uint           `json:"crawler_run_id"`
	CrawlerRun       *CrawlerRun     `json:"-"                  gorm:"foreignKey:CrawlerRunID"`
	GeoDataID        *uint           `json:"geo_data_id"`
	GeoData          *GeoData        `json:"geo_data"           gorm:"foreignKey:GeoDataID"`
	IndustryID       *uint           `json:"industry_id"`
	Industry         *Industry       `json:"industry"           gorm:"foreignKey:IndustryID"`
	OccupationID     *uint           `json:"occupation_id"`
	Occupation       *Occupation     `json:"occupation"         gorm:"foreignKey:OccupationID"`
	SourceID         *uint           `json:"source_id"`
	Source           *Source         `json:"source"             gorm:"foreignKey:SourceID"`
	ExperienceID     *uint           `json:"experience_id"`
	Experience       *Experience     `json:"experience"         gorm:"foreignKey:ExperienceID"`
	PostedAt         *time.Time      `json:"posted_at"`
	EndDate          *time.Time      `json:"end_date"`
	MinhashSignature pq.Int64Array   `json:"minhash_signature"  gorm:"type:integer[]"`
	ConfidenceScore  float64         `json:"confidence_score"   gorm:"type:decimal(5,4)"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (JobMetaData) TableName() string { return "meta_data" }

// ── LSH dedup index ───────────────────────────────────────────────────────────

type LshIndex struct {
	BucketKey string   `json:"bucket_key" gorm:"primaryKey;size:64;not null"`
	BandNo    int      `json:"band_no"    gorm:"not null"`
	JobPostID uint     `json:"job_post_id" gorm:"primaryKey;not null"`
	JobPost   *JobPost `json:"-"          gorm:"foreignKey:JobPostID;constraint:OnDelete:CASCADE;"`
}

func (LshIndex) TableName() string { return "lsh_index" }

// ── API payloads ──────────────────────────────────────────────────────────────

type BatchSavePayload struct {
	NewJobs    []JobPost  `json:"new_jobs"    binding:"required"`
	LshIndexes []LshIndex `json:"lsh_indexes" binding:"required"`
}

type BatchUpdatePayload struct {
	Duplicates []JobMetaData `json:"duplicates" binding:"required"`
}

type ReconciliationPayload struct {
	CrawlerRunID uint `json:"crawler_run_id" binding:"required"`
}

type SkillDemand struct {
	ID           uint   `json:"id"`
	Skill        string `json:"skill"`
	OpenJobCount int64  `json:"open_job_count"`
}

type EmployerDemand struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}

type CrawlTimeGap struct {
	LastCrawledAt  *time.Time `json:"last_crawled_at"`
	GapSeconds     float64    `json:"gap_seconds"`
	GapHuman       string     `json:"gap_human"`
}

type SourceJobCount struct {
	ID           uint   `json:"id"`
	Source       string `json:"source"`
	OpenJobCount int64  `json:"open_job_count"`
}

type JobCountWithTrend struct {
	ActiveJobCount  int64   `json:"active_job_count"`
	LastMonthCount  int64   `json:"last_month_count"`
	ChangePercent   float64 `json:"change_percent"`
	Trend           string  `json:"trend"` // "up", "down", "stable"
}

type OccupationJobCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}

type IndustryJobCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}

type ExperienceJobCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}

type EducationLevelJobCount struct {
	ID           uint   `json:"id"`
	Level        string `json:"level"`
	OpenJobCount int64  `json:"open_job_count"`
}

type RemoteOnSiteCount struct {
	RemoteCount  int64 `json:"remote_count"`
	OnSiteCount  int64 `json:"on_site_count"`
}

type JobTypeJobCount struct {
	ID           uint   `json:"id"`
	Type         string `json:"type"`
	OpenJobCount int64  `json:"open_job_count"`
}

type OccupationYearlyTrend struct {
	Year         int   `json:"year"`
	OpenJobCount int64 `json:"open_job_count"`
}

type TopJobRole struct {
	JobRole      string `json:"job_role"`
	OpenJobCount int64  `json:"open_job_count"`
}

type IndustryYearlyTrend struct {
	Year         int   `json:"year"`
	OpenJobCount int64 `json:"open_job_count"`
}

type ExperienceYearlyJobCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}

type ProvinceJobCount struct {
	ID           uint    `json:"id"`
	Province     string  `json:"province"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lng"`
	OpenJobCount int64   `json:"open_job_count"`
}

type EducationLevelYearlyJobCount struct {
	ID           uint   `json:"id"`
	Level        string `json:"level"`
	OpenJobCount int64  `json:"open_job_count"`
}

type TopEmployerByIndustryYear struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OpenJobCount int64  `json:"open_job_count"`
}