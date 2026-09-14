package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"marketlens-go-backend/repositories"
)


func registerStatsTools(server *mcp.Server, repo *repositories.JobRepository) {
	registerNoArgTool(server, "get_active_job_stats",
		"Get the current active job count, with month-over-month trend.",
		func() (any, error) { return repo.GetActiveJobCountWithTrend() })

	registerNoArgTool(server, "get_stats_by_occupation",
		"Get current active job counts grouped by occupation major group.",
		func() (any, error) { return repo.GetActiveJobCountByOccupation() })

	registerNoArgTool(server, "get_stats_by_industry",
		"Get current active job counts grouped by industry sector.",
		func() (any, error) { return repo.GetActiveJobCountByIndustry() })

	registerNoArgTool(server, "get_stats_by_experience",
		"Get current active job counts grouped by experience level.",
		func() (any, error) { return repo.GetActiveJobCountByExperience() })


	registerNoArgTool(server, "get_stats_by_formality",
		"Get current active job counts grouped by formal/informal sector.",
		func() (any, error) { return repo.GetActiveJobCountByFormality() })

	registerNoArgTool(server, "get_stats_by_employment_sector",
		"Get current active job counts grouped by employment sector.",
		func() (any, error) { return repo.GetActiveJobCountByEmploymentSector() })

	registerNoArgTool(server, "get_stats_by_gender",
		"Get current active job counts grouped by gender.",
		func() (any, error) { return repo.GetActiveJobCountByGender() })

	registerNoArgTool(server, "get_stats_by_vocational_education",
		"Get current active job counts grouped by vocational education (NVQ) level.",
		func() (any, error) { return repo.GetActiveJobCountByVocationalEducation() })

	registerNoArgTool(server, "get_remote_vs_onsite",
		"Get current active job counts split between remote and on-site.",
		func() (any, error) { return repo.GetRemoteVsOnSiteCount() })

	registerNoArgTool(server, "get_stats_by_job_type",
		"Get current active job counts grouped by job type (Full Time, Part Time, Contract, Internship).",
		func() (any, error) { return repo.GetActiveJobCountByJobType() })
}