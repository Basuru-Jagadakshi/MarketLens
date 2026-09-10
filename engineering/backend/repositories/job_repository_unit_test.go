package repositories_test

import (
	"errors"
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	err = db.AutoMigrate(
		&models.JobPost{},
		&models.JobMetaData{},
		&models.JobType{},
		&models.Skill{},
		&models.GeoData{},
		&models.Employer{},
		&models.JobPostSkill{},
		&models.MajorGroup{},
		&models.SubMajorGroup{},
		&models.MinorGroup{},
		&models.UnitGroup{},
		&models.OccupationGroup{},
		&models.IndustrySector{},
		&models.IndustryDivision{},
		&models.IndustryGroup{},
		&models.IndustryClass{},
		&models.IndustrySubclass{},
		&models.EducationLevel{},
		&models.Gender{},
		&models.Formality{},
		&models.EmploymentSector{},
		&models.Experience{},
		&models.VocationalEducation{},
	)
	if err != nil {
		t.Fatalf("Schema generation auto-migration failed: %v", err)
	}

	db.Create(&models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612})
	db.Create(&models.GeoData{Province: "Central", Latitude: 7.2906, Longitude: 80.6337})

	repo := repositories.NewJobRepository(db)
	return db, repo
}

// ---------------------------------------------------------------------------
// GetAllMajorGroupsForDateRange - existence/deletion boundary tests.
// Same pattern applies to every other ...ForDateRange sibling
// (GetAllIndustrySectorsForDateRange, GetSubMajorGroupsByMajorGroupForDateRange,
// etc.) - only the model/table differs.
// ---------------------------------------------------------------------------

func TestGetAllMajorGroupsForDateRange_ExcludesEntityCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	// Created "today" - well after the range being queried.
	db.Create(&models.MajorGroup{Name: "Brand New Category", Code: "99"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1) // range ends yesterday

	results, err := repo.GetAllMajorGroupsForDateRange(fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Category", r.Name,
			"a major group created after the range's end date should not appear in a historical report")
	}
}

func TestGetAllMajorGroupsForDateRange_IncludesEntityDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "Retiring Category", Code: "50"}
	db.Create(&mg)

	// Soft-delete it "today".
	db.Delete(&mg)

	// Query a range that spans from well before the deletion through today -
	// the category was genuinely active for most of this window.
	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetAllMajorGroupsForDateRange(fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Retiring Category" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetAllMajorGroupsForDateRange_ExcludesEntityDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "Long Gone Category", Code: "51"}
	db.Create(&mg)
	db.Delete(&mg) // deleted "today"

	// Query a range that starts tomorrow - entirely after the deletion.
	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetAllMajorGroupsForDateRange(fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Category", r.Name,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetAllMajorGroupsForDateRange_IncludesEntityDeletedOnTheSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "Same Day Deletion", Code: "52"}
	db.Create(&mg)
	db.Delete(&mg) // deleted right now, same calendar day as fromDate below

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetAllMajorGroupsForDateRange(fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "an entity deleted on the same day the range starts should still appear for that day")
}

func TestGetAllMajorGroups_NeverAppliesDateFiltering(t *testing.T) {
	db, repo := setupTestDB(t)

	// The crawler-facing method should just be "currently active", regardless
	// of any date logic - soft-deleted rows must never appear here.
	active := models.MajorGroup{Name: "Active Category", Code: "1"}
	db.Create(&active)

	deleted := models.MajorGroup{Name: "Deleted Category", Code: "2"}
	db.Create(&deleted)
	db.Delete(&deleted)

	results, err := repo.GetAllMajorGroups()
	require.NoError(t, err)

	names := make([]string, 0, len(results))
	for _, r := range results {
		names = append(names, r.Name)
	}
	assert.Contains(t, names, "Active Category")
	assert.NotContains(t, names, "Deleted Category")
}

// ---------------------------------------------------------------------------
// GetLevelChildren - existence/deletion boundary tests at one level of the
// occupation hierarchy (major-group -> sub-major-group). Same pattern
// applies to every other level and to the industry branch.
// ---------------------------------------------------------------------------

func TestGetLevelChildren_ExcludesChildCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "Parent", Code: "1"}
	db.Create(&mg)

	db.Create(&models.SubMajorGroup{MajorGroupID: mg.ID, Name: "New Child", Code: "11"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	children, _, err := repo.GetLevelChildren("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, c := range children {
		assert.NotEqual(t, "New Child", c.Name)
	}
}

func TestGetLevelChildren_IncludesChildDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "Parent", Code: "1"}
	db.Create(&mg)

	child := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "Retiring Child", Code: "11"}
	db.Create(&child)
	db.Delete(&child)

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	children, childLevel, err := repo.GetLevelChildren("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, "sub-major-group", childLevel)

	found := false
	for _, c := range children {
		if c.Name == "Retiring Child" {
			found = true
		}
	}
	assert.True(t, found)
}

func TestGetLevelChildren_LeafLevelReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, _, err := repo.GetLevelChildren("occupation", "occupation-group", 1, time.Now(), time.Now())
	assert.Error(t, err, "leaf levels have no children and should return an error")
}

func TestGetLevelChildren_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, _, err := repo.GetLevelChildren("not-a-real-standard", "major-group", 1, time.Now(), time.Now())
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// GetGenderByLevel - representative of the seven ...ByLevel breakdown
// methods (GetFormalityByLevel, GetExperienceByLevel, etc. all share this
// exact created_at/deleted_at boundary pattern).
// ---------------------------------------------------------------------------

func buildMinimalOccupationHierarchy(t *testing.T, db interface {
	Create(interface{}) *gorm_DB
}) {
	t.Helper()
}

func TestGetGenderByLevel_ExcludesGenderCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	// Minimal occupation hierarchy chain down to occupation-group, since
	// buildJobPostIDsForLevel needs a real chain to resolve job post ids.
	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	// A brand-new gender category, created "today".
	db.Create(&models.Gender{GenderType: "Newly Added"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Newly Added", r.GenderType,
			"a gender category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetGenderByLevel_IncludesGenderDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	g := models.Gender{GenderType: "Retiring Gender"}
	db.Create(&g)
	db.Delete(&g)

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.GenderType == "Retiring Gender" {
			found = true
		}
	}
	assert.True(t, found, "a gender deleted mid-range should still appear in the breakdown for that range")
}

type occupationHierarchyFixture struct {
	MajorGroup      models.MajorGroup
	SubMajorGroup   models.SubMajorGroup
	MinorGroup      models.MinorGroup
	UnitGroup       models.UnitGroup
	OccupationGroup models.OccupationGroup
}

func buildOccupationHierarchyFixture(t *testing.T, db interface {
	Create(value interface{}) interface{ Error() error }
}) occupationHierarchyFixture {
	t.Helper()
	// (placeholder signature removed below - see actual helper)
	return occupationHierarchyFixture{}
}

func TestGetTopHiringEmployersByOccupationLevel_RanksEmployersByVacancyCountDescending(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	empA := models.Employer{Name: "Employer A"}
	db.Create(&empA)
	empB := models.Employer{Name: "Employer B"}
	db.Create(&empB)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Employer A: two postings, 3 + 2 = 5 total vacancies
	jobA1 := models.JobPost{EmployerID: empA.ID, JobTypeID: jobType.ID, JobRole: "Role A1", NoOfVacancies: 3}
	db.Create(&jobA1)
	db.Create(&models.JobMetaData{JobPostID: jobA1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	jobA2 := models.JobPost{EmployerID: empA.ID, JobTypeID: jobType.ID, JobRole: "Role A2", NoOfVacancies: 2}
	db.Create(&jobA2)
	db.Create(&models.JobMetaData{JobPostID: jobA2.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	// Employer B: one posting, 10 vacancies - should rank first
	jobB1 := models.JobPost{EmployerID: empB.ID, JobTypeID: jobType.ID, JobRole: "Role B1", NoOfVacancies: 10}
	db.Create(&jobB1)
	db.Create(&models.JobMetaData{JobPostID: jobB1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTopHiringEmployersByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.Equal(t, "Employer B", results[0].Name)
	assert.Equal(t, int64(10), results[0].OpenJobCount)
	assert.Equal(t, "Employer A", results[1].Name)
	assert.Equal(t, int64(5), results[1].OpenJobCount)
}

func TestGetTopHiringEmployersByOccupationLevel_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	emp := models.Employer{Name: "Outside Range Employer"}
	db.Create(&emp)

	// Posted well before the queried range.
	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{EmployerID: emp.ID, JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: oldPostedAt})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetTopHiringEmployersByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results, "a job posted well outside the queried range should not contribute to the results")
}

func TestGetTopHiringEmployersByOccupationLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	// Target hierarchy - the one we'll actually query.
	mgTarget := models.MajorGroup{Name: "Target MG", Code: "1"}
	db.Create(&mgTarget)
	smgTarget := models.SubMajorGroup{MajorGroupID: mgTarget.ID, Name: "SMG", Code: "11"}
	db.Create(&smgTarget)
	mingTarget := models.MinorGroup{SubMajorGroupID: smgTarget.ID, Name: "MinG", Code: "111"}
	db.Create(&mingTarget)
	ugTarget := models.UnitGroup{MinorGroupID: mingTarget.ID, Name: "UG", Code: "1111"}
	db.Create(&ugTarget)
	ogTarget := models.OccupationGroup{UnitGroupID: ugTarget.ID, Name: "OG", Code: "11111"}
	db.Create(&ogTarget)

	// A separate, unrelated hierarchy - jobs here must never appear in the
	// target major group's results.
	mgOther := models.MajorGroup{Name: "Other MG", Code: "2"}
	db.Create(&mgOther)
	smgOther := models.SubMajorGroup{MajorGroupID: mgOther.ID, Name: "SMG Other", Code: "21"}
	db.Create(&smgOther)
	mingOther := models.MinorGroup{SubMajorGroupID: smgOther.ID, Name: "MinG Other", Code: "211"}
	db.Create(&mingOther)
	ugOther := models.UnitGroup{MinorGroupID: mingOther.ID, Name: "UG Other", Code: "2111"}
	db.Create(&ugOther)
	ogOther := models.OccupationGroup{UnitGroupID: ugOther.ID, Name: "OG Other", Code: "21111"}
	db.Create(&ogOther)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	emp := models.Employer{Name: "Wrong Hierarchy Employer"}
	db.Create(&emp)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{EmployerID: emp.ID, JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 8}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTopHiringEmployersByOccupationLevel("major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results, "a job posted under a different major group must not appear in this major group's results")
}

func TestGetTopHiringEmployersByOccupationLevel_LimitsToTop5(t *testing.T) {
	db, repo := setupTestDB(t)

	mg := models.MajorGroup{Name: "MG", Code: "1"}
	db.Create(&mg)
	smg := models.SubMajorGroup{MajorGroupID: mg.ID, Name: "SMG", Code: "11"}
	db.Create(&smg)
	ming := models.MinorGroup{SubMajorGroupID: smg.ID, Name: "MinG", Code: "111"}
	db.Create(&ming)
	ug := models.UnitGroup{MinorGroupID: ming.ID, Name: "UG", Code: "1111"}
	db.Create(&ug)
	og := models.OccupationGroup{UnitGroupID: ug.ID, Name: "OG", Code: "11111"}
	db.Create(&og)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Seven distinct employers, each with a single posting - more than the
	// method's LIMIT 5.
	for i := 0; i < 7; i++ {
		emp := models.Employer{Name: "Employer " + string(rune('A'+i))}
		db.Create(&emp)

		job := models.JobPost{EmployerID: emp.ID, JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: uint(i + 1)}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTopHiringEmployersByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Len(t, results, 5, "results should be capped at 5 employers even when more exist")
}

func TestGetTopHiringEmployersByOccupationLevel_InvalidLevelReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetTopHiringEmployersByOccupationLevel("not-a-real-level", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err, "an invalid occupation level should surface the error from buildJobPostIDsForLevel")
}

func TestGetAllSkillsByOccupationLevel_ReturnsCorrectTotalIndependentOfLimit(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Three distinct skills, each attached to its own job post.
	skillNames := []string{"Go", "PostgreSQL", "Docker"}
	for _, name := range skillNames {
		skill := models.Skill{Skill: name}
		db.Create(&skill)

		job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 1}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
		db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	// Request only 2 rows back, but total should still report all 3.
	results, total, err := repo.GetAllSkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate, 2, 0)
	require.NoError(t, err)

	assert.Len(t, results, 2, "limit should cap the returned rows")
	assert.Equal(t, int64(3), total, "total must reflect the full matching set, unaffected by limit")
}

func TestGetAllSkillsByOccupationLevel_OffsetSkipsCorrectNumberOfRows(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Skills with distinct, ordered vacancy counts so ranking is deterministic:
	// "Highest" (10) > "Middle" (5) > "Lowest" (1)
	skillCounts := map[string]uint{"Highest": 10, "Middle": 5, "Lowest": 1}
	for name, count := range skillCounts {
		skill := models.Skill{Skill: name}
		db.Create(&skill)

		job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: count}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
		db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	// Page 1: first result only.
	page1, total, err := repo.GetAllSkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate, 1, 0)
	require.NoError(t, err)
	require.Len(t, page1, 1)
	assert.Equal(t, "Highest", page1[0].Skill)
	assert.Equal(t, int64(3), total)

	// Page 2: skip the first, get the second-ranked skill.
	page2, total, err := repo.GetAllSkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate, 1, 1)
	require.NoError(t, err)
	require.Len(t, page2, 1)
	assert.Equal(t, "Middle", page2[0].Skill)
	assert.Equal(t, int64(3), total, "total should be identical across pages of the same query")
}

func TestGetAllSkillsByOccupationLevel_SumsVacanciesAcrossMultipleJobsForSameSkill(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	skill := models.Skill{Skill: "Go"}
	db.Create(&skill)

	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job1.ID, SkillID: skill.ID})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job2.ID, SkillID: skill.ID})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, total, err := repo.GetAllSkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate, 10, 0)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(1), total, "one distinct skill, even though it's referenced by two job posts")
	assert.Equal(t, int64(5), results[0].OpenJobCount, "vacancy counts across all matching jobs for this skill should be summed")
}

func TestGetAllSkillsByOccupationLevel_ZeroLimitReturnsAllRows(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	for _, name := range []string{"Go", "Python", "Rust"} {
		skill := models.Skill{Skill: name}
		db.Create(&skill)
		job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 1}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
		db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	// limit=0 means "no limit applied" per the method's own `if limit > 0` guard.
	results, total, err := repo.GetAllSkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate, 0, 0)
	require.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, int64(3), total)
}

func TestGetAllSkillsByOccupationLevel_InvalidLevelReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, _, err := repo.GetAllSkillsByOccupationLevel("not-a-real-level", 1, time.Now().AddDate(0, 0, -10), time.Now(), 10, 0)
	assert.Error(t, err)
}

func TestGetTop15SkillsByOccupationLevel_OrdersByVacancyCountDescending(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	skillCounts := map[string]uint{"Low": 1, "High": 20, "Mid": 5}
	for name, count := range skillCounts {
		skill := models.Skill{Skill: name}
		db.Create(&skill)

		job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: count}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
		db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 3)

	assert.Equal(t, "High", results[0].Skill)
	assert.Equal(t, int64(20), results[0].OpenJobCount)
	assert.Equal(t, "Mid", results[1].Skill)
	assert.Equal(t, int64(5), results[1].OpenJobCount)
	assert.Equal(t, "Low", results[2].Skill)
	assert.Equal(t, int64(1), results[2].OpenJobCount)
}

func TestGetTop15SkillsByOccupationLevel_CapsAt15EvenWithMoreMatches(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	// 18 distinct skills - more than the method's LIMIT 15.
	for i := 0; i < 18; i++ {
		skill := models.Skill{Skill: "Skill-" + string(rune('A'+i))}
		db.Create(&skill)

		job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: uint(i + 1)}
		db.Create(&job)
		db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
		db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})
	}

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Len(t, results, 15, "results should be capped at 15 skills even when more exist")
}

func TestGetTop15SkillsByOccupationLevel_SumsVacanciesAcrossMultipleJobsForSameSkill(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	skill := models.Skill{Skill: "Go"}
	db.Create(&skill)

	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 4}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job1.ID, SkillID: skill.ID})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 6}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job2.ID, SkillID: skill.ID})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(10), results[0].OpenJobCount, "vacancy counts across all matching jobs for this skill should be summed")
}

func TestGetTop15SkillsByOccupationLevel_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	skill := models.Skill{Skill: "Outdated Skill"}
	db.Create(&skill)

	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 9}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: oldPostedAt})
	db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results, "a job posted well outside the queried range should not contribute to the results")
}

func TestGetTop15SkillsByOccupationLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, ogTarget := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	skill := models.Skill{Skill: "Wrong Hierarchy Skill"}
	db.Create(&skill)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 8}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, PostedAt: postedAt})
	db.Create(&models.JobPostSkill{JobPostID: job.ID, SkillID: skill.ID})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results, "a job posted under a different major group must not appear in this major group's results")

	// Sanity check: the "other" hierarchy's own occupation group id is real
	// and distinct, confirming the two fixtures didn't accidentally collide.
	assert.NotEqual(t, ogTarget.ID, ogOther.ID)
}

func TestGetTop15SkillsByOccupationLevel_NoMatchingJobsReturnsEmptySlice(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetTop15SkillsByOccupationLevel("major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetTop15SkillsByOccupationLevel_InvalidLevelReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetTop15SkillsByOccupationLevel("not-a-real-level", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetJobTypeByLevel_SumsVacanciesPerJobType(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	fullTime := models.JobType{Type: "Full Time"}
	db.Create(&fullTime)
	partTime := models.JobType{Type: "Part Time"}
	db.Create(&partTime)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Two Full Time postings (3 + 2 = 5), one Part Time posting (1).
	job1 := models.JobPost{JobTypeID: fullTime.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: fullTime.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: partTime.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetJobTypeByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byType := make(map[string]int64)
	for _, r := range results {
		byType[r.Type] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byType["Full Time"])
	assert.Equal(t, int64(1), byType["Part Time"])
}

func TestGetJobTypeByLevel_IncludesJobTypesWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	usedType := models.JobType{Type: "Full Time"}
	db.Create(&usedType)
	unusedType := models.JobType{Type: "Internship"}
	db.Create(&unusedType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: usedType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetJobTypeByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "every job type should appear, even ones with no matching jobs, since the joins are LEFT JOINs")

	byType := make(map[string]int64)
	for _, r := range results {
		byType[r.Type] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byType["Full Time"])
	assert.Equal(t, int64(0), byType["Internship"], "a job type with zero matching jobs should show a count of 0, not be omitted")
}

func TestGetJobTypeByLevel_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 9}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: oldPostedAt})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetJobTypeByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the job type row still appears, per the LEFT JOIN behavior")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted outside the range should not contribute to the count")
}

func TestGetJobTypeByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, ogTarget := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetJobTypeByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")

	assert.NotEqual(t, ogTarget.ID, ogOther.ID)
}

func TestGetJobTypeByLevel_WorksForIndustryStandardToo(t *testing.T) {
	db, repo := setupTestDB(t)

	industrySector := models.IndustrySector{Name: "Sector", Code: "1"}
	db.Create(&industrySector)
	industryDivision := models.IndustryDivision{IndustrySectorID: industrySector.ID, Name: "Division", Code: "11"}
	db.Create(&industryDivision)
	industryGroup := models.IndustryGroup{IndustryDivisionID: industryDivision.ID, Name: "Group", Code: "111"}
	db.Create(&industryGroup)
	industryClass := models.IndustryClass{IndustryGroupID: industryGroup.ID, Name: "Class", Code: "1111"}
	db.Create(&industryClass)
	industrySubclass := models.IndustrySubclass{IndustryClassID: industryClass.ID, Name: "Subclass", Code: "11111"}
	db.Create(&industrySubclass)

	jobType := models.JobType{Type: "Contract"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 3}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, IndustrySubclassID: industrySubclass.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetJobTypeByLevel("industry", "industry-sector", industrySector.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "Contract", results[0].Type)
	assert.Equal(t, int64(3), results[0].OpenJobCount)
}

func TestGetJobTypeByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetJobTypeByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetRemoteOnSiteByLevel_SplitsCountsByIsRemoteFlag(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)

	// Two remote postings (5 + 3 = 8), one on-site posting (2).
	remoteJob1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Remote 1", NoOfVacancies: 5, IsRemote: true}
	db.Create(&remoteJob1)
	db.Create(&models.JobMetaData{JobPostID: remoteJob1.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	remoteJob2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Remote 2", NoOfVacancies: 3, IsRemote: true}
	db.Create(&remoteJob2)
	db.Create(&models.JobMetaData{JobPostID: remoteJob2.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	onSiteJob := models.JobPost{JobTypeID: jobType.ID, JobRole: "On-Site", NoOfVacancies: 2, IsRemote: false}
	db.Create(&onSiteJob)
	db.Create(&models.JobMetaData{JobPostID: onSiteJob.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	result, err := repo.GetRemoteOnSiteByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(8), result.RemoteCount)
	assert.Equal(t, int64(2), result.OnSiteCount)
}

func TestGetRemoteOnSiteByLevel_NoMatchingJobsReturnsZeroForBoth(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	result, err := repo.GetRemoteOnSiteByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.RemoteCount, "with no matching jobs at all, RemoteCount should default to its zero value, not be left unset in a way that panics or errors")
	assert.Equal(t, int64(0), result.OnSiteCount)
}

func TestGetRemoteOnSiteByLevel_OnlyRemoteJobsLeavesOnSiteAtZero(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Remote Only", NoOfVacancies: 4, IsRemote: true}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	result, err := repo.GetRemoteOnSiteByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(4), result.RemoteCount)
	assert.Equal(t, int64(0), result.OnSiteCount, "on-site count should default to 0 when no on-site rows exist in the grouped result at all")
}

func TestGetRemoteOnSiteByLevel_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 9, IsRemote: true}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: oldPostedAt})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	result, err := repo.GetRemoteOnSiteByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.RemoteCount, "a job posted well outside the queried range should not contribute to the counts")
	assert.Equal(t, int64(0), result.OnSiteCount)
}

func TestGetRemoteOnSiteByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 6, IsRemote: true}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	result, err := repo.GetRemoteOnSiteByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.RemoteCount, "a job posted under a different major group must not count toward this one")
	assert.Equal(t, int64(0), result.OnSiteCount)
}

func TestGetRemoteOnSiteByLevel_InvalidLevelReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetRemoteOnSiteByLevel("occupation", "not-a-real-level", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetVocationalEducationByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	nvq4 := models.VocationalEducation{Level: "NVQ 4"}
	db.Create(&nvq4)
	nvq5 := models.VocationalEducation{Level: "NVQ 5"}
	db.Create(&nvq5)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, VocationalEducationID: nvq4.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, VocationalEducationID: nvq4.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, VocationalEducationID: nvq5.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byLevel := make(map[string]int64)
	for _, r := range results {
		byLevel[r.Level] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byLevel["NVQ 4"])
	assert.Equal(t, int64(1), byLevel["NVQ 5"])
}

func TestGetVocationalEducationByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.VocationalEducation{Level: "NVQ 4"}
	db.Create(&used)
	unused := models.VocationalEducation{Level: "NVQ 7"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, VocationalEducationID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	byLevel := make(map[string]int64)
	for _, r := range results {
		byLevel[r.Level] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byLevel["NVQ 4"])
	assert.Equal(t, int64(0), byLevel["NVQ 7"])
}

func TestGetVocationalEducationByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	// Created "today" - after the range being queried.
	db.Create(&models.VocationalEducation{Level: "Brand New NVQ"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New NVQ", r.Level,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetVocationalEducationByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	nvq := models.VocationalEducation{Level: "Retiring NVQ"}
	db.Create(&nvq)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, VocationalEducationID: nvq.ID, PostedAt: postedAt})

	db.Delete(&nvq) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Level == "Retiring NVQ" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetVocationalEducationByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	nvq := models.VocationalEducation{Level: "Long Gone NVQ"}
	db.Create(&nvq)
	db.Delete(&nvq) // deleted "today"

	// Query a range that starts tomorrow - entirely after the deletion.
	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone NVQ", r.Level,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetVocationalEducationByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	nvq := models.VocationalEducation{Level: "Same Day Deletion"}
	db.Create(&nvq)
	db.Delete(&nvq) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Level == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetVocationalEducationByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	nvq := models.VocationalEducation{Level: "NVQ 4"}
	db.Create(&nvq)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, VocationalEducationID: nvq.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetVocationalEducationByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the NVQ category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetVocationalEducationByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetVocationalEducationByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetGenderByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	male := models.Gender{GenderType: "Male"}
	db.Create(&male)
	female := models.Gender{GenderType: "Female"}
	db.Create(&female)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, GenderID: male.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, GenderID: male.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, GenderID: female.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byGender := make(map[string]int64)
	for _, r := range results {
		byGender[r.GenderType] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byGender["Male"])
	assert.Equal(t, int64(1), byGender["Female"])
}

func TestGetGenderByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.Gender{GenderType: "Male"}
	db.Create(&used)
	unused := models.Gender{GenderType: "Female"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GenderID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	byGender := make(map[string]int64)
	for _, r := range results {
		byGender[r.GenderType] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byGender["Male"])
	assert.Equal(t, int64(0), byGender["Female"])
}

func TestGetGenderByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	db.Create(&models.Gender{GenderType: "Brand New Gender"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Gender", r.GenderType,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetGenderByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	gender := models.Gender{GenderType: "Retiring Gender"}
	db.Create(&gender)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GenderID: gender.ID, PostedAt: postedAt})

	db.Delete(&gender) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.GenderType == "Retiring Gender" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetGenderByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	gender := models.Gender{GenderType: "Long Gone Gender"}
	db.Create(&gender)
	db.Delete(&gender) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Gender", r.GenderType,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetGenderByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	gender := models.Gender{GenderType: "Same Day Deletion"}
	db.Create(&gender)
	db.Delete(&gender) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetGenderByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.GenderType == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetGenderByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	gender := models.Gender{GenderType: "Male"}
	db.Create(&gender)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, GenderID: gender.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetGenderByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the gender category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetGenderByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetGenderByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetFormalityByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	formal := models.Formality{FormalityType: "Formal"}
	db.Create(&formal)
	informal := models.Formality{FormalityType: "Informal"}
	db.Create(&informal)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, FormalityID: formal.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, FormalityID: formal.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, FormalityID: informal.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byType := make(map[string]int64)
	for _, r := range results {
		byType[r.FormalityType] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byType["Formal"])
	assert.Equal(t, int64(1), byType["Informal"])
}

func TestGetFormalityByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.Formality{FormalityType: "Formal"}
	db.Create(&used)
	unused := models.Formality{FormalityType: "Informal"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, FormalityID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	byType := make(map[string]int64)
	for _, r := range results {
		byType[r.FormalityType] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byType["Formal"])
	assert.Equal(t, int64(0), byType["Informal"])
}

func TestGetFormalityByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	db.Create(&models.Formality{FormalityType: "Brand New Formality"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Formality", r.FormalityType,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetFormalityByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	formality := models.Formality{FormalityType: "Retiring Formality"}
	db.Create(&formality)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, FormalityID: formality.ID, PostedAt: postedAt})

	db.Delete(&formality) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.FormalityType == "Retiring Formality" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetFormalityByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	formality := models.Formality{FormalityType: "Long Gone Formality"}
	db.Create(&formality)
	db.Delete(&formality) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Formality", r.FormalityType,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetFormalityByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	formality := models.Formality{FormalityType: "Same Day Deletion"}
	db.Create(&formality)
	db.Delete(&formality) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.FormalityType == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetFormalityByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	formality := models.Formality{FormalityType: "Formal"}
	db.Create(&formality)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, FormalityID: formality.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetFormalityByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the formality category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetFormalityByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetFormalityByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetEducationLevelByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	degree := models.EducationLevel{Level: "Degree"}
	db.Create(&degree)
	alevel := models.EducationLevel{Level: "A/L"}
	db.Create(&alevel)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, EducationLevelID: degree.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, EducationLevelID: degree.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, EducationLevelID: alevel.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byLevel := make(map[string]int64)
	for _, r := range results {
		byLevel[r.Level] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byLevel["Degree"])
	assert.Equal(t, int64(1), byLevel["A/L"])
}

func TestGetEducationLevelByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.EducationLevel{Level: "Degree"}
	db.Create(&used)
	unused := models.EducationLevel{Level: "O/L"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EducationLevelID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	byLevel := make(map[string]int64)
	for _, r := range results {
		byLevel[r.Level] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byLevel["Degree"])
	assert.Equal(t, int64(0), byLevel["O/L"])
}

func TestGetEducationLevelByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	db.Create(&models.EducationLevel{Level: "Brand New Level"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Level", r.Level,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetEducationLevelByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	edu := models.EducationLevel{Level: "Retiring Level"}
	db.Create(&edu)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EducationLevelID: edu.ID, PostedAt: postedAt})

	db.Delete(&edu) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Level == "Retiring Level" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetEducationLevelByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	edu := models.EducationLevel{Level: "Long Gone Level"}
	db.Create(&edu)
	db.Delete(&edu) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Level", r.Level,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetEducationLevelByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	edu := models.EducationLevel{Level: "Same Day Deletion"}
	db.Create(&edu)
	db.Delete(&edu) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Level == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetEducationLevelByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	edu := models.EducationLevel{Level: "Degree"}
	db.Create(&edu)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, EducationLevelID: edu.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEducationLevelByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the education level category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetEducationLevelByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetEducationLevelByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetProvinceByLevel_SumsVacanciesPerProvince(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	western := models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612}
	db.Create(&western)
	central := models.GeoData{Province: "Central", Latitude: 7.2906, Longitude: 80.6337}
	db.Create(&central)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, GeoDataID: western.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, GeoDataID: western.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, GeoDataID: central.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetProvinceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byProvince := make(map[string]int64)
	for _, r := range results {
		byProvince[r.Province] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byProvince["Western"])
	assert.Equal(t, int64(1), byProvince["Central"])
}

func TestGetProvinceByLevel_IncludesProvinceWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612}
	db.Create(&used)
	unused := models.GeoData{Province: "Southern", Latitude: 6.0535, Longitude: 80.2210}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GeoDataID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetProvinceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a province with zero matching jobs should still appear, per the LEFT JOIN")

	byProvince := make(map[string]int64)
	for _, r := range results {
		byProvince[r.Province] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byProvince["Western"])
	assert.Equal(t, int64(0), byProvince["Southern"])
}

func TestGetProvinceByLevel_ReturnsLatitudeAndLongitude(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	western := models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612}
	db.Create(&western)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 2}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GeoDataID: western.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetProvinceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.InDelta(t, 6.9271, results[0].Latitude, 0.0001)
	assert.InDelta(t, 79.8612, results[0].Longitude, 0.0001)
}

func TestGetProvinceByLevel_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	province := models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612}
	db.Create(&province)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 9}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, GeoDataID: province.ID, PostedAt: oldPostedAt})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetProvinceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the province row still appears, per the LEFT JOIN behavior")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted outside the range should not contribute to the count")
}

func TestGetProvinceByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	province := models.GeoData{Province: "Western", Latitude: 6.9271, Longitude: 79.8612}
	db.Create(&province)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, GeoDataID: province.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetProvinceByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the province row itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetProvinceByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetProvinceByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetExperienceByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	entry := models.Experience{Name: "Entry Level"}
	db.Create(&entry)
	senior := models.Experience{Name: "Senior"}
	db.Create(&senior)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, ExperienceID: entry.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, ExperienceID: entry.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, ExperienceID: senior.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byName := make(map[string]int64)
	for _, r := range results {
		byName[r.Name] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byName["Entry Level"])
	assert.Equal(t, int64(1), byName["Senior"])
}

func TestGetExperienceByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.Experience{Name: "Entry Level"}
	db.Create(&used)
	unused := models.Experience{Name: "Senior"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, ExperienceID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	byName := make(map[string]int64)
	for _, r := range results {
		byName[r.Name] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byName["Entry Level"])
	assert.Equal(t, int64(0), byName["Senior"])
}

func TestGetExperienceByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	db.Create(&models.Experience{Name: "Brand New Experience"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Experience", r.Name,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetExperienceByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	exp := models.Experience{Name: "Retiring Experience"}
	db.Create(&exp)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, ExperienceID: exp.ID, PostedAt: postedAt})

	db.Delete(&exp) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Retiring Experience" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetExperienceByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	exp := models.Experience{Name: "Long Gone Experience"}
	db.Create(&exp)
	db.Delete(&exp) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Experience", r.Name,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetExperienceByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	exp := models.Experience{Name: "Same Day Deletion"}
	db.Create(&exp)
	db.Delete(&exp) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetExperienceByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	exp := models.Experience{Name: "Entry Level"}
	db.Create(&exp)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, ExperienceID: exp.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetExperienceByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the experience category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetExperienceByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetExperienceByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetEmploymentSectorByLevel_SumsVacanciesPerCategory(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	private := models.EmploymentSector{Sector: "Private"}
	db.Create(&private)
	government := models.EmploymentSector{Sector: "Government"}
	db.Create(&government)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: og.ID, EmploymentSectorID: private.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: og.ID, EmploymentSectorID: private.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: og.ID, EmploymentSectorID: government.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	bySector := make(map[string]int64)
	for _, r := range results {
		bySector[r.Sector] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), bySector["Private"])
	assert.Equal(t, int64(1), bySector["Government"])
}

func TestGetEmploymentSectorByLevel_IncludesCategoryWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	used := models.EmploymentSector{Sector: "Private"}
	db.Create(&used)
	unused := models.EmploymentSector{Sector: "NGO"}
	db.Create(&unused)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EmploymentSectorID: used.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a category with zero matching jobs should still appear, per the LEFT JOIN")

	bySector := make(map[string]int64)
	for _, r := range results {
		bySector[r.Sector] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), bySector["Private"])
	assert.Equal(t, int64(0), bySector["NGO"])
}

func TestGetEmploymentSectorByLevel_ExcludesCategoryCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	db.Create(&models.EmploymentSector{Sector: "Brand New Sector"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Sector", r.Sector,
			"a category created after the range's end date should not appear in a historical breakdown")
	}
}

func TestGetEmploymentSectorByLevel_IncludesCategoryDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "MG", "1")

	sector := models.EmploymentSector{Sector: "Retiring Sector"}
	db.Create(&sector)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, EmploymentSectorID: sector.ID, PostedAt: postedAt})

	db.Delete(&sector) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Sector == "Retiring Sector" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a category deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetEmploymentSectorByLevel_ExcludesCategoryDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	sector := models.EmploymentSector{Sector: "Long Gone Sector"}
	db.Create(&sector)
	db.Delete(&sector) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Sector", r.Sector,
			"a category already deleted before the range started should not appear")
	}
}

func TestGetEmploymentSectorByLevel_IncludesCategoryDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "MG", "1")

	sector := models.EmploymentSector{Sector: "Same Day Deletion"}
	db.Create(&sector)
	db.Delete(&sector) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mg.ID, fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Sector == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a category deleted on the same day the range starts should still appear for that day")
}

func TestGetEmploymentSectorByLevel_ExcludesJobsUnderDifferentMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgTarget, _ := seedOccupationHierarchy(t, db, "Target MG", "1")
	_, ogOther := seedOccupationHierarchy(t, db, "Other MG", "2")

	sector := models.EmploymentSector{Sector: "Private"}
	db.Create(&sector)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Other Role", NoOfVacancies: 7}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: ogOther.ID, EmploymentSectorID: sector.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetEmploymentSectorByLevel("occupation", "major-group", mgTarget.ID, fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the employment sector category itself still appears since it exists in the table")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted under a different major group must not count toward this one")
}

func TestGetEmploymentSectorByLevel_InvalidStandardReturnsError(t *testing.T) {
	_, repo := setupTestDB(t)

	_, err := repo.GetEmploymentSectorByLevel("not-a-real-standard", "major-group", 1, time.Now().AddDate(0, 0, -10), time.Now())
	assert.Error(t, err)
}

func TestGetOccupationJobCountByDateRange_SumsVacanciesPerMajorGroup(t *testing.T) {
	db, repo := setupTestDB(t)

	mgA, ogA := seedOccupationHierarchy(t, db, "Group A", "1")
	mgB, ogB := seedOccupationHierarchy(t, db, "Group B", "2")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, OccupationGroupID: ogA.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 2}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, OccupationGroupID: ogA.ID, PostedAt: postedAt})

	job3 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 3", NoOfVacancies: 1}
	db.Create(&job3)
	db.Create(&models.JobMetaData{JobPostID: job3.ID, OccupationGroupID: ogB.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2)

	byName := make(map[string]int64)
	for _, r := range results {
		byName[r.Name] = r.OpenJobCount
	}
	assert.Equal(t, int64(5), byName[mgA.Name])
	assert.Equal(t, int64(1), byName[mgB.Name])
}

func TestGetOccupationJobCountByDateRange_IncludesGroupWithZeroMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	usedMg, usedOg := seedOccupationHierarchy(t, db, "Used Group", "1")
	unusedMg, _ := seedOccupationHierarchy(t, db, "Unused Group", "2")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 4}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: usedOg.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 2, "a major group with zero matching jobs should still appear, per the LEFT JOIN")

	byName := make(map[string]int64)
	for _, r := range results {
		byName[r.Name] = r.OpenJobCount
	}
	assert.Equal(t, int64(4), byName[usedMg.Name])
	assert.Equal(t, int64(0), byName[unusedMg.Name])
}

func TestGetOccupationJobCountByDateRange_ExcludesGroupCreatedAfterRange(t *testing.T) {
	db, repo := setupTestDB(t)

	_, _ = seedOccupationHierarchy(t, db, "Existing Group", "1")
	db.Create(&models.MajorGroup{Name: "Brand New Group", Code: "99"})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now().AddDate(0, 0, -1)

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Brand New Group", r.Name,
			"a major group created after the range's end date should not appear in a historical report")
	}
}

func TestGetOccupationJobCountByDateRange_IncludesGroupDeletedDuringRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "Retiring Group", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role", NoOfVacancies: 6}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: postedAt})

	db.Delete(&mg) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Retiring Group" {
			found = true
			assert.Equal(t, int64(6), r.OpenJobCount, "the historical count should still reflect the real postings made before deletion")
		}
	}
	assert.True(t, found, "a major group deleted mid-range should still appear, since it was active for part of the range")
}

func TestGetOccupationJobCountByDateRange_ExcludesGroupDeletedBeforeRangeStarted(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "Long Gone Group", "1")
	db.Delete(&mg) // deleted "today"

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 10)

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)

	for _, r := range results {
		assert.NotEqual(t, "Long Gone Group", r.Name,
			"a major group already deleted before the range started should not appear")
	}
}

func TestGetOccupationJobCountByDateRange_IncludesGroupDeletedOnSameDayAsFromDate(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, _ := seedOccupationHierarchy(t, db, "Same Day Deletion", "1")
	db.Delete(&mg) // deleted right now

	fromDate := time.Now().Truncate(24 * time.Hour) // midnight today
	toDate := time.Now().AddDate(0, 0, 5)

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)

	found := false
	for _, r := range results {
		if r.Name == "Same Day Deletion" {
			found = true
		}
	}
	assert.True(t, found, "a major group deleted on the same day the range starts should still appear for that day")
}

func TestGetOccupationJobCountByDateRange_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	mg, og := seedOccupationHierarchy(t, db, "Group", "1")

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	oldPostedAt := time.Now().AddDate(0, 0, -100)
	job := models.JobPost{JobTypeID: jobType.ID, JobRole: "Old Role", NoOfVacancies: 9}
	db.Create(&job)
	db.Create(&models.JobMetaData{JobPostID: job.ID, OccupationGroupID: og.ID, PostedAt: oldPostedAt})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)
	require.Len(t, results, 1, "the major group row still appears, per the LEFT JOIN behavior")
	assert.Equal(t, int64(0), results[0].OpenJobCount, "a job posted outside the range should not contribute to the count")
	assert.Equal(t, mg.Name, results[0].Name)
}

func TestGetOccupationJobCountByDateRange_NoMajorGroupsReturnsEmptySlice(t *testing.T) {
	_, repo := setupTestDB(t)

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	results, err := repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetTotalVacancyCount_SumsVacanciesAcrossAllMatchingJobs(t *testing.T) {
	db, repo := setupTestDB(t)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)
	postedAt := time.Now().AddDate(0, 0, -5)

	job1 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 1", NoOfVacancies: 3}
	db.Create(&job1)
	db.Create(&models.JobMetaData{JobPostID: job1.ID, PostedAt: postedAt})

	job2 := models.JobPost{JobTypeID: jobType.ID, JobRole: "Role 2", NoOfVacancies: 7}
	db.Create(&job2)
	db.Create(&models.JobMetaData{JobPostID: job2.ID, PostedAt: postedAt})

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	total, err := repo.GetTotalVacancyCount(fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestGetTotalVacancyCount_ExcludesJobsOutsideDateRange(t *testing.T) {
	db, repo := setupTestDB(t)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	insideRange := time.Now().AddDate(0, 0, -5)
	outsideRange := time.Now().AddDate(0, 0, -100)

	jobInside := models.JobPost{JobTypeID: jobType.ID, JobRole: "Inside", NoOfVacancies: 4}
	db.Create(&jobInside)
	db.Create(&models.JobMetaData{JobPostID: jobInside.ID, PostedAt: insideRange})

	jobOutside := models.JobPost{JobTypeID: jobType.ID, JobRole: "Outside", NoOfVacancies: 100}
	db.Create(&jobOutside)
	db.Create(&models.JobMetaData{JobPostID: jobOutside.ID, PostedAt: outsideRange})

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	total, err := repo.GetTotalVacancyCount(fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(4), total, "a job posted well outside the range should not contribute to the total")
}

func TestGetTotalVacancyCount_IncludesJobsOnBothBoundaryDates(t *testing.T) {
	db, repo := setupTestDB(t)

	jobType := models.JobType{Type: "Full Time"}
	db.Create(&jobType)

	fromDate := time.Now().AddDate(0, 0, -10)
	toDate := time.Now()

	// One job posted exactly on fromDate, one exactly on toDate - both
	// boundaries are inclusive per the BETWEEN clause.
	jobOnFrom := models.JobPost{JobTypeID: jobType.ID, JobRole: "On From", NoOfVacancies: 2}
	db.Create(&jobOnFrom)
	db.Create(&models.JobMetaData{JobPostID: jobOnFrom.ID, PostedAt: fromDate})

	jobOnTo := models.JobPost{JobTypeID: jobType.ID, JobRole: "On To", NoOfVacancies: 3}
	db.Create(&jobOnTo)
	db.Create(&models.JobMetaData{JobPostID: jobOnTo.ID, PostedAt: toDate})

	total, err := repo.GetTotalVacancyCount(fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total, "jobs posted exactly on either boundary date should be included")
}

func TestGetTotalVacancyCount_NoMatchingJobsReturnsZero(t *testing.T) {
	_, repo := setupTestDB(t)

	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	total, err := repo.GetTotalVacancyCount(fromDate, toDate)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total, "with no matching jobs at all, the total should be 0, not an error or a null-derived panic")
}