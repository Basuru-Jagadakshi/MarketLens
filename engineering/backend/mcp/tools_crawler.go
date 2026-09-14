package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"marketlens-go-backend/repositories"
)

func registerCrawlerTools(server *mcp.Server, repo *repositories.JobRepository) {
	registerNoArgTool(server, "get_crawler_last_job_count",
		"Get the job count from the most recently completed crawler run.",
		func() (any, error) { return repo.GetLastCrawledJobCount() })

	registerNoArgTool(server, "get_crawler_time_gap",
		"Get the time elapsed since the last completed crawler run.",
		func() (any, error) { return repo.GetTimeSinceLastCrawl() })

	registerNoArgTool(server, "get_crawler_runs",
		"List the most recent crawler runs with their status.",
		func() (any, error) { return repo.GetAllCrawlerRuns() })

	registerNoArgTool(server, "get_sources",
		"Get the sources with their crawled active job count.",
		func() (any, error) {return repo.GetSourcesWithActiveJobCount() })
}