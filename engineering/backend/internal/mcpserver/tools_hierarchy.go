package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerHierarchyTools(server *mcp.Server, svcs *services) {
	registerNoArgTool(server, "get_major_groups", "List all occupation major groups (top level of the SLSO occupation hierarchy).",
		svcs.majorGroups.GetAll)
	registerIDTool(server, "get_major_group", "Get a single occupation major group by id.",
		svcs.majorGroups.GetByID)

	registerNoArgTool(server, "get_sub_major_groups", "List all occupation sub major groups.",
		svcs.subMajorGroups.GetAll)
	registerIDTool(server, "get_sub_major_group", "Get a single occupation sub major group by id.",
		svcs.subMajorGroups.GetByID)

	registerNoArgTool(server, "get_minor_groups", "List all occupation minor groups.",
		svcs.minorGroups.GetAll)
	registerIDTool(server, "get_minor_group", "Get a single occupation minor group by id.",
		svcs.minorGroups.GetByID)

	registerNoArgTool(server, "get_unit_groups", "List all occupation unit groups.",
		svcs.unitGroups.GetAll)
	registerIDTool(server, "get_unit_group", "Get a single occupation unit group by id.",
		svcs.unitGroups.GetByID)

	// Paginated
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_occupation_groups",
		Description: "List occupation groups (leaf level of the SLSO hierarchy), paginated.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Limit  int `json:"limit" jsonschema:"max rows to return, default 20"`
		Offset int `json:"offset" jsonschema:"rows to skip, default 0"`
	}) (*mcp.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		items, total, err := svcs.occupationGroups.GetAllPaginated(ctx, limit, in.Offset)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"items": items, "total": total, "limit": limit, "offset": in.Offset}, nil
	})
	registerIDTool(server, "get_occupation_group", "Get a single occupation group by id.",
		svcs.occupationGroups.GetByID)

	registerNoArgTool(server, "get_industry_sectors", "List all industry sectors (top level of the SLSIC industry hierarchy).",
		svcs.industrySectors.GetAll)
	registerIDTool(server, "get_industry_sector", "Get a single industry sector by id.",
		svcs.industrySectors.GetByID)

	registerNoArgTool(server, "get_industry_divisions", "List all industry divisions.",
		svcs.industryDivisions.GetAll)
	registerIDTool(server, "get_industry_division", "Get a single industry division by id.",
		svcs.industryDivisions.GetByID)

	registerNoArgTool(server, "get_industry_groups", "List all industry groups.",
		svcs.industryGroups.GetAll)
	registerIDTool(server, "get_industry_group", "Get a single industry group by id.",
		svcs.industryGroups.GetByID)

	registerNoArgTool(server, "get_industry_classes", "List all industry classes.",
		svcs.industryClasses.GetAll)
	registerIDTool(server, "get_industry_class", "Get a single industry class by id.",
		svcs.industryClasses.GetByID)

	// Paginated
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_industry_subclasses",
		Description: "List industry subclasses (leaf level of the SLSIC hierarchy), paginated.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Limit  int `json:"limit" jsonschema:"max rows to return, default 20"`
		Offset int `json:"offset" jsonschema:"rows to skip, default 0"`
	}) (*mcp.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		items, total, err := svcs.industrySubclasses.GetAllPaginated(ctx, limit, in.Offset)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"items": items, "total": total, "limit": limit, "offset": in.Offset}, nil
	})
	registerIDTool(server, "get_industry_subclass", "Get a single industry subclass by id.",
		svcs.industrySubclasses.GetByID)

	registerChildrenTool(server, svcs)
	registerLevelBreakdownTool(server, "get_total_job_count_by_level",
		"Get the total vacancy count for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetTotalJobCountByLevel)
	registerLevelBreakdownTool(server, "get_employment_sector_by_level",
		"Get employment sector breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetEmploymentSectorByLevel)
	registerLevelBreakdownTool(server, "get_experience_by_level",
		"Get experience-level breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetExperienceByLevel)
	registerLevelBreakdownTool(server, "get_province_by_level",
		"Get province breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetProvinceByLevel)
	registerLevelBreakdownTool(server, "get_education_by_level",
		"Get education level breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetEducationLevelByLevel)
	registerLevelBreakdownTool(server, "get_formality_by_level",
		"Get formal/informal breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetFormalityByLevel)
	registerLevelBreakdownTool(server, "get_gender_by_level",
		"Get gender breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetGenderByLevel)
	registerLevelBreakdownTool(server, "get_vocational_education_by_level",
		"Get vocational education (NVQ) breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetVocationalEducationByLevel)
	registerLevelBreakdownTool(server, "get_remote_onsite_by_level",
		"Get remote vs on-site job counts for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetRemoteOnSiteHybridByLevel)
	registerLevelBreakdownTool(server, "get_job_type_by_level",
		"Get job type breakdown for a given occupation/industry hierarchy level and id, within a date range.",
		svcs.analytics.GetJobTypeByLevel)
}

// levelInput mirrors the /:standard/:level/:id/... path params, plus the
// from/to date query params every breakdown endpoint requires.
type levelInput struct {
	Standard string `json:"standard" jsonschema:"'occupation' or 'industry'"`
	Level    string `json:"level" jsonschema:"the hierarchy level name, e.g. 'major-group' or 'industry-sector'"`
	ID       uint   `json:"id" jsonschema:"the numeric id at that level"`
	FromDate string `json:"from_date" jsonschema:"start date, format YYYY-MM-DD"`
	ToDate   string `json:"to_date" jsonschema:"end date, format YYYY-MM-DD"`
}

// parseLevelInput parses the shared level/id/date fields used by every
// generic breakdown tool. Semantic checks (standard/level validity, date
// order) are the service layer's job, since they're business rules shared
// with the HTTP handlers, not a parsing concern.
func parseLevelInput(in levelInput) (uint, time.Time, time.Time, error) {
	from, to, err := parseDateRange(in.FromDate, in.ToDate)
	if err != nil {
		return 0, time.Time{}, time.Time{}, err
	}
	return in.ID, from, to, nil
}

// registerLevelBreakdownTool registers one of the ten generic
// standard/level/id breakdown tools, given the service call it should
// dispatch to.
func registerLevelBreakdownTool[T any](
	server *mcp.Server,
	name, description string,
	fn func(ctx context.Context, standard, level string, id uint, from, to time.Time) (T, error),
) {
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: description},
		func(ctx context.Context, _ *mcp.CallToolRequest, in levelInput) (*mcp.CallToolResult, any, error) {
			id, from, to, err := parseLevelInput(in)
			if err != nil {
				return nil, nil, err
			}
			out, err := fn(ctx, in.Standard, in.Level, id, from, to)
			if err != nil {
				return nil, nil, err
			}
			return nil, wrapIfList(out), nil
		},
	)
}

// registerChildrenTool registers the "get the next level down" tool
// (mirrors GET /:standard/:level/:id/children).
func registerChildrenTool(server *mcp.Server, svcs *services) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "get_hierarchy_children",
		Description: "Get the immediate child-level entities under a given occupation/industry hierarchy " +
			"level and id, each with its aggregated job count for a date range. Fails for leaf levels " +
			"('occupation-group', 'industry-subclass'), which have no children.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in levelInput) (*mcp.CallToolResult, any, error) {
		id, from, to, err := parseLevelInput(in)
		if err != nil {
			return nil, nil, err
		}
		items, childLevel, err := svcs.analytics.GetLevelChildren(ctx, in.Standard, in.Level, id, from, to)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]any{"child_level": childLevel, "children": items}, nil
	})
}
