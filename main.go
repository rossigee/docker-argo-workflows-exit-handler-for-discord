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
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Color       int                    `json:"color"`
	Fields      []map[string]string    `json:"fields"`
}

func sendDiscordMessage(webhookURL string, embeds []DiscordEmbed) {
	payload := map[string][]DiscordEmbed{"embeds": embeds}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
        	os.Exit(1)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending request:", err)
        	os.Exit(1)
	}

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Failed to send message, Status:", resp.Status)
        	os.Exit(1)
	}
}

func main() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	workflowStatus := os.Getenv("ARGO_WORKFLOW_STATUS")

	if webhookURL == "" || workflowStatus == "" {
		fmt.Println("Required environment variables missing.")
		return
	}

	var nodes []NodeInfo
	input := []byte(os.Getenv("ARGO_FAILED_NODES"))
	json.Unmarshal(input, &nodes)

	var embeds []DiscordEmbed

	for _, node := range nodes {
		embed := DiscordEmbed{
			Title:       "Node Failure Information",
			Description: fmt.Sprintf("Node: %s", node.DisplayName),
			Color:       16711680,  // Bright red
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

	sendDiscordMessage(webhookURL, embeds)
}
