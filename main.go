package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NodeInfo struct {
	DisplayName  string `json:"displayName"`
	Message      string `json:"message"`
	TemplateName string `json:"templateName"`
	Phase        string `json:"phase"`
	PodName      string `json:"podName"`
	FinishedAt   string `json:"finishedAt"`
}

type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []map[string]string `json:"fields"`
}

func prepareDiscordMessage(embeds []DiscordEmbed) ([]byte, error) {
	payload := map[string][]DiscordEmbed{"embeds": embeds}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %s", err)
	}

	payloadFile := os.Getenv("DISCORD_PAYLOAD_TO_FILE")
	if payloadFile != "" {
		file, err := os.Create(payloadFile)
		if err != nil {
			fmt.Printf("error creating payload file: %s", err)
		}
		_, err = file.Write(jsonPayload)
		if err != nil {
			fmt.Printf("error writing to payload file: %s", err)
		}
	}

	return jsonPayload, nil
}

func sendDiscordMessage(webhookURL string, jsonPayload []byte) error {
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error sending request: %s", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message, Status: %s", resp.Status)
	}

	return nil
}

func main() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	workflowStatus := os.Getenv("ARGO_WORKFLOW_STATUS")

	if webhookURL == "" || workflowStatus == "" {
		fmt.Println("Required environment variables missing.")
		os.Exit(1)
	}

	var nodes []NodeInfo
	input := []byte(os.Getenv("ARGO_FAILED_NODES"))
	json.Unmarshal(input, &nodes)

	var embeds []DiscordEmbed

	for _, node := range nodes {
		embed := DiscordEmbed{
			Title:       "Node Failure Information",
			Description: fmt.Sprintf("Node: %s", node.DisplayName),
			Color:       16711680, // Bright red
			Fields: []map[string]string{
				{"name": "Message", "value": node.Message},
				{"name": "Template", "value": node.TemplateName},
				{"name": "Phase", "value": node.Phase},
				{"name": "Pod Name", "value": node.PodName},
				{"name": "Finished At", "value": node.FinishedAt},
			},
		}
		embeds = append(embeds, embed)
	}

	payload, err := prepareDiscordMessage(embeds)
	if err != nil {
		fmt.Printf("Error preparing Discord message: %s", err)
		os.Exit(1)
	}
	err = sendDiscordMessage(webhookURL, payload)
	if err != nil {
		fmt.Printf("Error sending Discord message: %s", err)
		os.Exit(1)
	}
}
