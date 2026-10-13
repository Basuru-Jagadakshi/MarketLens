package mcpserver

import (
	"context"
	"net/http"
	"encoding/json"
	"os"          
	"strings"
	"reflect"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"marketlens-go-backend/internal/service"
)

// services bundles every service the MCP tools call — the same
// service.CRUDService/AnalyticsService instances the HTTP handlers use, so
// validation and business rules live in exactly one place regardless of
// which transport a request comes in on.
type services struct {
	analytics *service.AnalyticsService

	experiences          *service.CRUDService[models.Experience]
	geoData              *service.CRUDService[models.GeoData]
	jobTypes             *service.CRUDService[models.JobType]
	employmentSectors    *service.CRUDService[models.EmploymentSector]
	educationLevels      *service.CRUDService[models.EducationLevel]
	formalities          *service.CRUDService[models.Formality]
	genders              *service.CRUDService[models.Gender]
	vocationalEducations *service.CRUDService[models.VocationalEducation]

	majorGroups      *service.CRUDService[models.MajorGroup]
	subMajorGroups   *service.CRUDService[models.SubMajorGroup]
	minorGroups      *service.CRUDService[models.MinorGroup]
	unitGroups       *service.CRUDService[models.UnitGroup]
	occupationGroups *service.CRUDService[models.OccupationGroup]

	industrySectors    *service.CRUDService[models.IndustrySector]
	industryDivisions  *service.CRUDService[models.IndustryDivision]
	industryGroups     *service.CRUDService[models.IndustryGroup]
	industryClasses    *service.CRUDService[models.IndustryClass]
	industrySubclasses *service.CRUDService[models.IndustrySubclass]
}

func newServices(store *repositories.Store) *services {
	return &services{
		analytics: service.NewAnalyticsService(store),

		experiences:          service.NewExperienceService(store),
		geoData:              service.NewGeoDataService(store),
		jobTypes:             service.NewJobTypeService(store),
		employmentSectors:    service.NewEmploymentSectorService(store),
		educationLevels:      service.NewEducationLevelService(store),
		formalities:          service.NewFormalityService(store),
		genders:              service.NewGenderService(store),
		vocationalEducations: service.NewVocationalEducationService(store),

		majorGroups:      service.NewMajorGroupService(store),
		subMajorGroups:   service.NewSubMajorGroupService(store),
		minorGroups:      service.NewMinorGroupService(store),
		unitGroups:       service.NewUnitGroupService(store),
		occupationGroups: service.NewOccupationGroupService(store),

		industrySectors:    service.NewIndustrySectorService(store),
		industryDivisions:  service.NewIndustryDivisionService(store),
		industryGroups:     service.NewIndustryGroupService(store),
		industryClasses:    service.NewIndustryClassService(store),
		industrySubclasses: service.NewIndustrySubclassService(store),
	}
}

const (
	protectedResPath = "/.well-known/oauth-protected-resource/mcp"
)

type contextKey string

func mustGetEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable: " + name)
	}
	return v
}

func mcpResourceID() string      { return mustGetEnv("MCP_RESOURCE_ID") }
func thunderIssuerURL() string   { return mustGetEnv("THUNDER_ISSUER") }
func mcpPublicBaseURL() string   { return mustGetEnv("MCP_PUBLIC_BASE_URL") }

func serveProtectedResourceMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"resource":                 mcpResourceID(),        
		"authorization_servers":    []string{thunderIssuerURL()}, 
		"scopes_supported": []string{"submit-newspaper-vacancies"},
		"bearer_methods_supported": []string{"header"},
	})
}

const (
	bearerTokenKey contextKey = "bearer_token"
	scopesKey      contextKey = "scopes"        
)

func requireValidToken(next http.Handler) http.Handler {   
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("WWW-Authenticate",
				`Bearer resource_metadata="`+mcpPublicBaseURL()+protectedResPath+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		_, scopes, err := verifyMCPToken(rawToken)  
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), bearerTokenKey, rawToken)
		ctx = context.WithValue(ctx, scopesKey, scopes)   
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// New builds the MCP server and registers every read-only (GET-equivalent)
// tool against the given repository's services.
func New(repo *repositories.Store) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "marketlens-mcp",
		Version: "v1.0.0",
	}, nil)

	svcs := newServices(repo)
	registerLookupTools(server, svcs)
	registerHierarchyTools(server, svcs)
	registerAnalysisTools(server, svcs)
	registerManualUploadTools(server)

	return server
}

// StartHTTP serves the given MCP server over Streamable HTTP at addr
// (e.g. ":9090"). 
func StartHTTP(server *mcp.Server, addr string) error {
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.HandleFunc(protectedResPath, serveProtectedResourceMetadata)
	mux.Handle("/mcp", requireValidToken(handler))   

	return http.ListenAndServe(addr, mux)
}

func wrapIfList(v any) any {
	if v == nil {
		return v
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return map[string]any{"items": v}
	case reflect.Map, reflect.Struct, reflect.Ptr:
		return v 
	default:
		return map[string]any{"value": v} 
	}
}

// emptyInput is used for tools that take no parameters at all.
type emptyInput struct{}

// registerNoArgTool registers a tool that takes no input and returns
// whatever the given repository call returns.
func registerNoArgTool[T any](server *mcp.Server, name, description string, fn func(ctx context.Context) (T, error)) {
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: description},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
			out, err := fn(ctx)
			if err != nil {
				return nil, nil, err
			}
			return nil, wrapIfList(out), nil
		},
	)
}

// idInput is used for tools scoped to a single numeric id (e.g. an
// industry_sector id, a major_group id).
type idInput struct {
	ID uint `json:"id" jsonschema:"the numeric id to look up"`
}

// registerIDTool registers a tool that takes a single "id" parameter.
func registerIDTool[T any](server *mcp.Server, name, description string, fn func(ctx context.Context, id uint) (T, error)) {
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: description},
		func(ctx context.Context, _ *mcp.CallToolRequest, in idInput) (*mcp.CallToolResult, any, error) {
			out, err := fn(ctx, in.ID)
			if err != nil {
				return nil, nil, err
			}
			return nil, wrapIfList(out), nil
		},
	)
}

// dateRangeInput is used for tools scoped to a from/to date window.
type dateRangeInput struct {
	FromDate string `json:"from_date" jsonschema:"start date, format YYYY-MM-DD"`
	ToDate   string `json:"to_date" jsonschema:"end date, format YYYY-MM-DD"`
}

// idAndYearInput is used for the occupation-scoped year-filtered tools
// (by-formality, by-gender, top-job-roles).
type idAndYearInput struct {
	ID   uint `json:"id" jsonschema:"the numeric id to look up"`
	Year int  `json:"year" jsonschema:"the calendar year, e.g. 2026"`
}

// idAndDateRangeInput is used for tools scoped to both a single id and a
// from/to date window (e.g. industry-scoped year analytics).
type idAndDateRangeInput struct {
	ID       uint   `json:"id" jsonschema:"the numeric id to look up"`
	FromDate string `json:"from_date" jsonschema:"start date, format YYYY-MM-DD"`
	ToDate   string `json:"to_date" jsonschema:"end date, format YYYY-MM-DD"`
}

func registerLookupTools(server *mcp.Server, svcs *services) {
	registerNoArgTool(server, "get_experiences", "List all experience levels.",
		svcs.experiences.GetAll)

	registerNoArgTool(server, "get_provinces", "List all Sri Lankan provinces used for geo-tagging job posts.",
		svcs.geoData.GetAll)

	registerNoArgTool(server, "get_job_types", "List all job types (Full Time, Part Time, Contract, Internship).",
		svcs.jobTypes.GetAll)

	registerNoArgTool(server, "get_employment_sectors", "List all employment sectors (Government, Private, NGO, etc.).",
		svcs.employmentSectors.GetAll)

	registerNoArgTool(server, "get_education_levels", "List all education level categories.",
		svcs.educationLevels.GetAll)

	registerNoArgTool(server, "get_formalities", "List all formality categories (Formal / Informal sector).",
		svcs.formalities.GetAll)

	registerNoArgTool(server, "get_genders", "List all gender categories used in job postings.",
		svcs.genders.GetAll)

	registerNoArgTool(server, "get_vocational_educations", "List all vocational education (NVQ) levels.",
		svcs.vocationalEducations.GetAll)
}