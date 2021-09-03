package config

import (
	"errors"
	"os"
)

type DiscordConfig struct {
	BotID             string
	TriggerEmojiAlias string
}

func GetDiscordConfig() (*DiscordConfig, error) {
	botID := os.Getenv("DISCORD_BOT_ID")
	triggerEmoji := os.Getenv("DISCORD_TRIGGER_EMOJI")
	if botID == "" || triggerEmoji == "" {
		return nil, errors.New("[Config] env is empty")
	}

	return &DiscordConfig{
		BotID:             botID,
		TriggerEmojiAlias: triggerEmoji,
	}, nil
}
