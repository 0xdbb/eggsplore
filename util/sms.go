package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ArkeselSMSPayload struct {
	Sender        string   `json:"sender"`
	Message       string   `json:"message"`
	Recipients    []string `json:"recipients"`
	ScheduledDate string   `json:"scheduled_date,omitempty"`
	CallbackURL   string   `json:"callback_url,omitempty"`
	UseCase       string   `json:"use_case,omitempty"`
	Sandbox       bool     `json:"sandbox,omitempty"`
}

type ArkeselSMSResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendSMS(
	recipient string,
	message string,
	senderID string,
	apiKey string,
	options ...map[string]any,
) error {
	payload := ArkeselSMSPayload{
		Sender:     senderID,
		Message:    message,
		Recipients: []string{recipient},
	}

	if len(options) > 0 {
		opt := options[0]
		if val, ok := opt["scheduled_date"].(string); ok {
			payload.ScheduledDate = val
		}
		if val, ok := opt["callback_url"].(string); ok {
			payload.CallbackURL = val
		}
		if val, ok := opt["use_case"].(string); ok {
			payload.UseCase = val
		}
		if val, ok := opt["sandbox"].(bool); ok {
			payload.Sandbox = val
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://sms.arkesel.com/api/v2/sms/send", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("SMS send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("arkesel returned status code: %d", resp.StatusCode)
	}

	var response ArkeselSMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse Arkesel response: %w", err)
	}
	log.Println("✅ Arkesel response:", response.Status, response.Message)

	if response.Status != "success" {
		return fmt.Errorf("arkesel error: %s - %s", response.Status, response.Message)
	}

	fmt.Printf("✅ Arkesel response: %+v\n", response)
	return nil
}
