package crawler_test

import (
	"context"
	"strings"
	"testing"

	"marketlens-go-backend/crawler"
	"marketlens-go-backend/crawler/mocks"
	"marketlens-go-backend/models"
)

// extractionResponder returns a full, valid extraction JSON for any
// prompt that looks like the extraction call (starts with the
// instruction template's opening line), and a classifier pick
// (mocks.FirstListedID) for anything else.
func extractionResponder(content string) mocks.Responder {
	return func(prompt string) string {
		if strings.HasPrefix(prompt, "Extract structural data from") {
			return content
		}
		id := mocks.FirstListedID(prompt)
		return `{"id": ` + itoa(id) + `}`
	}
}

func TestProcessBatch_DuplicateAndNewJobTogether(t *testing.T) {
	repo := mocks.NewMockRepository()
	llm := mocks.NewMockDeepSeek(t, extractionResponder(`{
		"job_type": {"type": "Full Time"},
		"work_mode": "onsite",
		"no_of_vacancies": 2,
		"meta_data": {
			"geo_data": {"province": "Western"},
			"posted_at": "2026-09-01T00:00:00Z",
			"confidence_score": 0.92,
			"formality_id": 1,
			"gender_id": 3,
			"vocational_education_id": 1,
			"employment_sector_id": 1,
			"education_level_id": 3,
			"experience_id": 2
		},
		"skills": [{"skill": "Go"}, {"skill": "PostgreSQL"}]
	}`))
	metaBuilder := crawler.NewMetadataBuilder(repo, 0)
	svc := crawler.NewIngestionService(repo, llm, metaBuilder, 0.65)

	firstJob := crawler.RawJobInput{
		Employer:     "Ceylon Software Solutions",
		JobRole:      "Senior Backend Engineer",
		Location:     "Colombo 03, Sri Lanka",
		Description:  "Go, PostgreSQL and Docker experience required for our growing backend team.",
		CrawlerRunID: 100,
		Source:       "Ikman",
	}

	// Seed the store with firstJob already saved, as if a previous
	// crawl run had inserted it.
	if err := svc.ProcessBatch(context.Background(), []crawler.RawJobInput{firstJob}); err != nil {
		t.Fatalf("seeding first batch failed: %v", err)
	}

	duplicateJob := crawler.RawJobInput{
		Employer:     "Ceylon Software Solutions",
		JobRole:      "Senior Backend Engineer",
		Location:     "Colombo 03",
		Description:  "Go, PostgreSQL, and Docker experience required for our growing backend team!",
		CrawlerRunID: 200,
		Source:       "Ikman",
	}
	newJob := crawler.RawJobInput{
		Employer:     "Lanka Retail Group",
		JobRole:      "Store Manager",
		Location:     "Kandy, Sri Lanka",
		Description:  "Manage daily operations, staff scheduling, and inventory for our flagship retail outlet.",
		CrawlerRunID: 200,
		Source:       "TopJobs",
	}

	if err := svc.ProcessBatch(context.Background(), []crawler.RawJobInput{duplicateJob, newJob}); err != nil {
		t.Fatalf("ProcessBatch failed: %v", err)
	}

	// The duplicate must trigger exactly one update, not a new save.
	if len(repo.UpdatedDuplicates) != 1 {
		t.Fatalf("expected 1 duplicate update, got %d", len(repo.UpdatedDuplicates))
	}
	if repo.UpdatedDuplicates[0].JobPostID != 1 {
		t.Errorf("expected duplicate update to target job ID 1, got %d", repo.UpdatedDuplicates[0].JobPostID)
	}
	if repo.UpdatedDuplicates[0].CrawlerRunID == nil || *repo.UpdatedDuplicates[0].CrawlerRunID != 200 {
		t.Errorf("expected duplicate's CrawlerRunID updated to 200, got %v", repo.UpdatedDuplicates[0].CrawlerRunID)
	}

	// The new job must trigger exactly one additional save (on top of
	// the first job saved during setup), fully extracted and
	// classified. Find it by location rather than assuming index 0,
	// since SavedJobs accumulates everything saved across both calls.
	if len(repo.SavedJobs) != 2 {
		t.Fatalf("expected 2 total jobs saved (1 setup + 1 new), got %d", len(repo.SavedJobs))
	}
	var saved *models.JobPost
	for i := range repo.SavedJobs {
		if repo.SavedJobs[i].Location == "Kandy, Sri Lanka" {
			saved = &repo.SavedJobs[i]
		}
	}
	if saved == nil {
		t.Fatal("new job (Kandy, Sri Lanka) not found among saved jobs")
	}
	if saved.JobType == nil || saved.JobType.Type != "Full Time" {
		t.Errorf("expected extracted job_type Full Time, got %+v", saved.JobType)
	}
	if len(saved.Skills) != 2 {
		t.Errorf("expected 2 extracted skills, got %d", len(saved.Skills))
	}
	if saved.MetaData.IndustrySubclassID == nil || *saved.MetaData.IndustrySubclassID == 0 {
		t.Errorf("expected a non-zero classified industry subclass id, got %v", saved.MetaData.IndustrySubclassID)
	}
	if saved.MetaData.OccupationGroupID == nil || *saved.MetaData.OccupationGroupID == 0 {
		t.Errorf("expected a non-zero classified occupation group id, got %v", saved.MetaData.OccupationGroupID)
	}
	if saved.MetaData.MinhashSignature == nil {
		t.Error("expected a stored MinHash signature on the new job")
	}
}

func TestProcessBatch_OmittedFieldsBecomeNilNotFKBreakingZero(t *testing.T) {
	repo := mocks.NewMockRepository()
	// Deliberately omit work_mode, formality_id, gender_id from the
	// extraction response, and make the classifier picks return 0 (no
	// match) by using an id that isn't in any offered option list.
	llm := mocks.NewMockDeepSeek(t, func(prompt string) string {
		if strings.HasPrefix(prompt, "Extract structural data from") {
			return `{
				"job_type": {"type": "Full Time"},
				"no_of_vacancies": 1,
				"meta_data": {
					"geo_data": {"province": "Western"},
					"posted_at": "2026-09-01T00:00:00Z",
					"confidence_score": 0.5
				},
				"skills": []
			}`
		}
		return `{"id": 0}`
	})
	metaBuilder := crawler.NewMetadataBuilder(repo, 0)
	svc := crawler.NewIngestionService(repo, llm, metaBuilder, 0.65)

	raw := crawler.RawJobInput{
		Employer:     "Test Co",
		JobRole:      "Tester",
		Location:     "Colombo",
		Description:  "Testing incomplete extraction handling",
		CrawlerRunID: 1,
		Source:       "Test",
	}

	if err := svc.ProcessBatch(context.Background(), []crawler.RawJobInput{raw}); err != nil {
		t.Fatalf("ProcessBatch failed: %v", err)
	}

	if len(repo.SavedJobs) != 1 {
		t.Fatalf("expected 1 job saved, got %d", len(repo.SavedJobs))
	}
	saved := repo.SavedJobs[0]

	if saved.MetaData.FormalityID != nil {
		t.Errorf("expected FormalityID nil (omitted by LLM), got pointer to %d", *saved.MetaData.FormalityID)
	}
	if saved.MetaData.GenderID != nil {
		t.Errorf("expected GenderID nil (omitted by LLM), got pointer to %d", *saved.MetaData.GenderID)
	}
	if saved.MetaData.IndustrySubclassID != nil {
		t.Errorf("expected IndustrySubclassID nil (no classification match), got pointer to %d", *saved.MetaData.IndustrySubclassID)
	}
	if saved.MetaData.OccupationGroupID != nil {
		t.Errorf("expected OccupationGroupID nil (no classification match), got pointer to %d", *saved.MetaData.OccupationGroupID)
	}
	if saved.WorkMode != "onsite" {
		t.Errorf("expected WorkMode to default to 'onsite' when omitted, got %q", saved.WorkMode)
	}
}

func TestProcessBatch_HandlesMarkdownFencedExtractionResponse(t *testing.T) {
	repo := mocks.NewMockRepository()
	fenced := "```json\n" + `{
		"job_type": {"type": "Part Time"},
		"work_mode": "remote",
		"no_of_vacancies": 1,
		"meta_data": {
			"geo_data": {"province": "Western"},
			"posted_at": "2026-09-01T00:00:00Z",
			"confidence_score": 0.8,
			"formality_id": 1,
			"gender_id": 3,
			"vocational_education_id": 1,
			"employment_sector_id": 1,
			"education_level_id": 3,
			"experience_id": 2
		},
		"skills": [{"skill": "Python"}]
	}` + "\n```"

	llm := mocks.NewMockDeepSeek(t, extractionResponder(fenced))
	metaBuilder := crawler.NewMetadataBuilder(repo, 0)
	svc := crawler.NewIngestionService(repo, llm, metaBuilder, 0.65)

	raw := crawler.RawJobInput{
		Employer:     "Fence Test Co",
		JobRole:      "Remote Analyst",
		Location:     "Colombo",
		Description:  "A job whose extraction response is wrapped in a markdown code fence.",
		CrawlerRunID: 1,
		Source:       "Test",
	}

	if err := svc.ProcessBatch(context.Background(), []crawler.RawJobInput{raw}); err != nil {
		t.Fatalf("ProcessBatch failed on fenced extraction response: %v", err)
	}

	if len(repo.SavedJobs) != 1 {
		t.Fatalf("expected 1 job saved, got %d", len(repo.SavedJobs))
	}
	if repo.SavedJobs[0].WorkMode != "remote" {
		t.Errorf("expected WorkMode 'remote' parsed correctly from fenced response, got %q", repo.SavedJobs[0].WorkMode)
	}
}