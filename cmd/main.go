package main

import (
	"log"

	"github.com/qbaware/repo-recap/internal/github"
	"github.com/qbaware/repo-recap/internal/messaging"
)

const (
	webhookURL = "https://discord.com/api/webhooks/1419089464759947316/HeOtlW-NgS6frv3aqUueu9VvckFJY1BW1tpqko9DJU5YF2yoC_4_IWXGoGNJuhxtltgH"
)

func main() {
	// Testing Discord messaging
	discordClient := messaging.NewDiscordClient(webhookURL)
	err := discordClient.SendWebhookMessage(messaging.DiscordWebhookMessage{
		Content:  "Repo recap completed!",
		Username: "Repo Bot",
		Embeds: []messaging.DiscordEmbed{{
			Title:       "Summary",
			Description: "5 commits, 12 files changed",
			Color:       5814783,
		}},
	})
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}

	// Testing GitHub API interaction
	ghClient := github.NewClient("your-token")
	ghClient.GetDailyCommitSummary("qbaware", "homeassistant-eldom")
}
