package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// ErrTransientExtractionFailure marks a failure in the extract/classify
// pipeline as retryable: a network problem, the LLM timing out (408) or
// rate limiting (429) or erroring (5xx), the LLM responding with
// unparseable/empty content, or a DB read failing while fetching
// reference data. None of these are about the job's own content — a
// retry later should succeed.
var ErrTransientExtractionFailure = errors.New("transient extraction failure")

// Option is the common (id, name) shape every hierarchy level's rows are
// converted into before being shown to the LLM.
type Option struct {
	ID   uint
	Name string
}

// toOptions converts a slice of any model into []Option, given a
// function that knows how to pull the id/name off that specific type.
func toOptions[T any](items []T, extract func(T) (uint, string)) []Option {
	options := make([]Option, len(items))
	for i, item := range items {
		id, name := extract(item)
		options[i] = Option{ID: id, Name: name}
	}
	return options
}

// LLMClient asks a chat-completion model either to pick one option from a
// list (AskToPick) or to answer a free-form prompt (Chat). The provider,
// endpoint, and model are configured entirely through environment
// variables — this package has no knowledge of which vendor is behind
// LLM_API_URL.
type LLMClient struct {
	httpClient *http.Client
	apiKey     string
	apiURL     string
	model      string
}

// NewLLMClient builds a client reading LLM_API_KEY, LLM_API_URL, and
// LLM_MODEL from the environment.
func NewLLMClient(httpClient *http.Client) *LLMClient {
	return &LLMClient{
		httpClient: httpClient,
		apiKey:     os.Getenv("LLM_API_KEY"),
		apiURL:     os.Getenv("LLM_API_URL"),
		model:      os.Getenv("LLM_MODEL"),
	}
}

// NewLLMClientWithURL is the same as NewLLMClient but lets the endpoint be
// overridden — useful for pointing at a mock server in tests.
func NewLLMClientWithURL(httpClient *http.Client, apiURL string) *LLMClient {
	c := NewLLMClient(httpClient)
	c.apiURL = apiURL
	return c
}

// Validate reports which required environment variables are missing.
// Call this once at startup so a misconfigured deployment fails fast
// instead of every extraction/classification call failing at request
// time with an opaque "unauthorized" or "connection refused" error.
func (c *LLMClient) Validate() error {
	var missing []string
	if c.apiKey == "" {
		missing = append(missing, "LLM_API_KEY")
	}
	if c.apiURL == "" {
		missing = append(missing, "LLM_API_URL")
	}
	if c.model == "" {
		missing = append(missing, "LLM_MODEL")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variable(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

type llmRequest struct {
	Model       string       `json:"model"`
	Messages    []llmMessage `json:"messages"`
	Temperature float64      `json:"temperature"`
	// ResponseFormat asks the model to return strict JSON with no
	// markdown code-fence wrapping. Requires the word "json" to appear
	// somewhere in the prompt, which every prompt built in this package
	// already satisfies ("Respond ONLY with JSON...").
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type pickResult struct {
	ID uint `json:"id"`
}

// Chat sends a single-turn prompt to the LLM and returns the raw response
// content string. AskToPick and job extraction both build on this instead
// of duplicating request/response handling.
func (c *LLMClient) Chat(ctx context.Context, prompt string) (string, error) {
	reqBody, err := json.Marshal(llmRequest{
		Model:          c.model,
		Messages:       []llmMessage{{Role: "user", Content: prompt}},
		Temperature:    0.0,
		ResponseFormat: &responseFormat{Type: "json_object"},
	})
	if err != nil {
		return "", fmt.Errorf("marshaling llm request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("building llm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: calling llm: %v", ErrTransientExtractionFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return "", fmt.Errorf("%w: llm returned status %d", ErrTransientExtractionFailure, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm returned status %d", resp.StatusCode)
	}

	var parsed llmResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("%w: decoding llm response: %v", ErrTransientExtractionFailure, err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("%w: llm response had no choices", ErrTransientExtractionFailure)
	}

	return parsed.Choices[0].Message.Content, nil
}

// stripJSONFences removes a leading/trailing markdown code fence
// (```json ... ``` or plain ``` ... ```) if present. Belt-and-suspenders
// alongside response_format: json_object — models occasionally wrap
// JSON in a fence anyway, and this costs nothing when there isn't one.
func stripJSONFences(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	if nl := strings.IndexByte(s, '\n'); nl != -1 && strings.TrimSpace(s[:nl]) != "" {
		// first line was a language tag like "json" — drop it
		s = s[nl+1:]
	}
	s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	return strings.TrimSpace(s)
}

// AskToPick asks the LLM to choose the single best matching option and
// returns its id (0 if there were no options to choose from).
func (c *LLMClient) AskToPick(ctx context.Context, jobText string, options []Option, levelName string) (uint, error) {
	if len(options) == 0 {
		return 0, nil
	}

	prompt := buildAskToPickPrompt(jobText, options, levelName)

	content, err := c.Chat(ctx, prompt)
	if err != nil {
		return 0, err
	}

	var result pickResult
	if err := json.Unmarshal([]byte(stripJSONFences(content)), &result); err != nil {
		return 0, fmt.Errorf("%w: parsing chosen id from llm content: %v", ErrTransientExtractionFailure, err)
	}

	for _, o := range options {
		if o.ID == result.ID {
			return result.ID, nil
		}
	}

	// The LLM returned an id that wasn't actually one of the options it
	// was given — a hallucination, not a valid pick. Treated the same
	// as "no match" (0) rather than trusted, since a fabricated nonzero
	// id would otherwise pass straight through as a plausible-looking
	// but invalid foreign key.
	log.Printf("llm returned id %d for %q, which was not among the %d offered options — treating as no match", result.ID, levelName, len(options))
	return 0, nil
}

// LevelFetcher fetches the option list for one hierarchy level, given
// the id chosen at the previous level.
type LevelFetcher func(parentID uint) ([]Option, error)

// Level is one step of a classification walk.
type Level struct {
	Name  string
	Fetch LevelFetcher
}

// Walk asks the LLM to pick one option at each level in order, feeding
// the chosen id into the next level's Fetch. Stops early and returns 0
// the moment a level has no options or the LLM doesn't pick one.
func Walk(ctx context.Context, llm *LLMClient, jobText string, levels []Level) (uint, error) {
	var parentID uint
	for _, level := range levels {
		options, err := level.Fetch(parentID)
		if err != nil {
			return 0, fmt.Errorf("%w: fetching %s options: %v", ErrTransientExtractionFailure, level.Name, err)
		}

		id, err := llm.AskToPick(ctx, jobText, options, level.Name)
		if err != nil {
			return 0, fmt.Errorf("asking LLM to pick %s: %w", level.Name, err)
		}
		if id == 0 {
			return 0, nil
		}
		parentID = id
	}
	return parentID, nil
}
