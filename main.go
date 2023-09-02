package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func sendDiscordMessage(webhookURL, message string) {
	payload := map[string]string{"content": message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Failed to send message, Status:", resp.Status)
	}
}

func main() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	workflowStatus := os.Getenv("ARGO_WORKFLOW_STATUS")

	if webhookURL == "" || workflowStatus == "" {
		fmt.Println("Required environment variables missing.")
		return
	}

	sendDiscordMessage(webhookURL, "Workflow finished with status: "+workflowStatus)
}

