package discord

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
)

const discordEmojiURL = "https://raw.githubusercontent.com/shokkunrf/discord-emoji/develop/nature.csv"

type EmojiSet struct {
	EmojiMap        map[string]string
	ReverseEmojiMap map[string]string
}

func GetEmojiSet() (*EmojiSet, error) {
	res, err := http.Get(discordEmojiURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	csvByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	csvAry, err := csv.NewReader(bytes.NewReader(csvByte)).ReadAll()
	if err != nil {
		return nil, err
	}

	emojiMap := make(map[string]string)
	reverseEmojiMap := make(map[string]string)
	for _, row := range csvAry[1:] {
		emojiMap[row[0]] = row[1]
		if reverseEmojiMap[row[1]] == "" {
			reverseEmojiMap[row[1]] = row[0]
		}
	}

	return &EmojiSet{
		EmojiMap:        emojiMap,
		ReverseEmojiMap: reverseEmojiMap,
	}, nil
}
