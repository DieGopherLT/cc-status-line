package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// StatusHook represents the JSON structure received from Claude Code's Status hook
type StatusHook struct {
	HookEventName   string         `json:"hook_event_name"`
	SessionID       string         `json:"session_id"`
	TranscriptPath  string         `json:"transcript_path"`
	CWD             string         `json:"cwd"`
	Model           Model          `json:"model"`
	Workspace       Workspace      `json:"workspace"`
	Version         string         `json:"version"`
	OutputStyle     Output         `json:"output_style"`
	Cost            Cost           `json:"cost"`
	ContextWindow   *ContextWindow `json:"context_window"`
}

// Model contains information about the Claude model being used
type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Workspace contains directory information
type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

// Output contains the output style configuration
type Output struct {
	Name string `json:"name"`
}

// Cost contains session cost and metrics
type Cost struct {
	TotalCostUSD        float64 `json:"total_cost_usd"`
	TotalDurationMS     int64   `json:"total_duration_ms"`
	TotalAPIDurationMS  int64   `json:"total_api_duration_ms"`
	TotalLinesAdded     int     `json:"total_lines_added"`
	TotalLinesRemoved   int     `json:"total_lines_removed"`
}

// CurrentUsage contains detailed token usage for the current request
type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// ContextWindow contains token usage information for the session
type ContextWindow struct {
	TotalInputTokens  int           `json:"total_input_tokens"`
	TotalOutputTokens int           `json:"total_output_tokens"`
	ContextWindowSize int           `json:"context_window_size"`
	CurrentUsage      *CurrentUsage `json:"current_usage"`
}

// ParseStatusHook reads and parses the status hook JSON from an io.Reader
func ParseStatusHook(reader io.Reader) (*StatusHook, error) {
	var hook StatusHook
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&hook); err != nil {
		return nil, fmt.Errorf("failed to parse status hook JSON: %w", err)
	}

	// Validate required fields
	if hook.Model.DisplayName == "" {
		return nil, fmt.Errorf("model.display_name is required")
	}
	if hook.Version == "" {
		return nil, fmt.Errorf("version is required")
	}

	return &hook, nil
}
