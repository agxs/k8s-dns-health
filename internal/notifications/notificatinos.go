package notifications

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

type NotificationType interface {
	Notify(err error)
}

type TeamsNotification struct{}

type RestartNotification struct{}

type CombinedNotification struct {}

func (n TeamsNotification) Notify(dnsError error) error {
	message := dnsError.Error()

	webhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
    if webhookURL == "" {
        return fmt.Errorf("environment variable TEAMS_WEBHOOK_URL is not set")
    }

    // Prepare the payload
    payload := map[string]string{
        "text": message,
    }

    // Encode the payload as JSON
    body, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    // Create a new HTTP request
    req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Set the appropriate header
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Check for a successful response (200 OK or 204 No Content)
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("received non-OK response: %d", resp.StatusCode)
    }

    return nil
}

func (n RestartNotification) Notify(err error) {
}

func (n CombinedNotification) Notify(err error) {
}
