package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordWebhookMessage represents a Discord webhook message
type DiscordWebhookMessage struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
	URL         string `json:"url,omitempty"`
}

// DiscordClient handles Discord API communication
type DiscordClient struct {
	WebhookURL string
}

// NewDiscordClient creates a new Discord client
func NewDiscordClient(webhookURL string) *DiscordClient {
	return &DiscordClient{
		WebhookURL: webhookURL,
	}
}

// SendWebhookMessage sends a message via Discord webhook
func (d *DiscordClient) SendWebhookMessage(message DiscordWebhookMessage) error {
	if d.WebhookURL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord API error: %s", resp.Status)
	}

	return nil
}
