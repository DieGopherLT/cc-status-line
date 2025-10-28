package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TranscriptLine represents a single line in the JSONL transcript file
type TranscriptLine struct {
	Timestamp          string   `json:"timestamp"`
	IsSidechain        bool     `json:"isSidechain"`
	IsAPIErrorMessage  bool     `json:"isApiErrorMessage"`
	Message            *Message `json:"message"`
}

// Message contains the assistant's response data
type Message struct {
	Usage *Usage `json:"usage"`
}

// Usage contains token usage information
type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
}

// TranscriptData contains parsed data from the transcript file
type TranscriptData struct {
	Lines          []TranscriptLine
	FirstTimestamp time.Time
	LastTimestamp  time.Time
}

// ParseTranscript reads and parses a JSONL transcript file
func ParseTranscript(filePath string) (*TranscriptData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open transcript file: %w", err)
	}
	defer file.Close()

	var lines []TranscriptLine
	var firstTimestamp, lastTimestamp time.Time

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var transcriptLine TranscriptLine
		if err := json.Unmarshal([]byte(line), &transcriptLine); err != nil {
			// Skip malformed lines
			continue
		}

		lines = append(lines, transcriptLine)

		// Track first and last timestamps
		if transcriptLine.Timestamp != "" {
			timestamp, err := time.Parse(time.RFC3339, transcriptLine.Timestamp)
			if err == nil {
				if firstTimestamp.IsZero() || timestamp.Before(firstTimestamp) {
					firstTimestamp = timestamp
				}
				if lastTimestamp.IsZero() || timestamp.After(lastTimestamp) {
					lastTimestamp = timestamp
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading transcript file: %w", err)
	}

	return &TranscriptData{
		Lines:          lines,
		FirstTimestamp: firstTimestamp,
		LastTimestamp:  lastTimestamp,
	}, nil
}
