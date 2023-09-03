package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSecondsToHumanReadable(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{86400, "1 days, 0 hours, 0 minutes, 0 seconds"},
		{3661, "0 days, 1 hours, 1 minutes, 1 seconds"},
	}

	for _, test := range tests {
		actual := secondsToHumanReadable(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, actual)
		}
	}
}

func TestPrepareDiscordMessage(t *testing.T) {
	embed := []DiscordEmbed{
		{
			Title:       "Test",
			Description: "Test Description",
			Color:       DISCORD_GREEN,
		},
	}

	expected := map[string][]DiscordEmbed{"embeds": embed}
	expectedJSON, _ := json.Marshal(expected)

	actual, err := prepareDiscordMessage(embed)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !bytes.Equal(actual, expectedJSON) {
		t.Errorf("Expected %s, got %s", string(expectedJSON), string(actual))
	}
}
