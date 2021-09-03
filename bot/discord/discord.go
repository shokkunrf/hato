package discord

import (
	"fmt"
	"hato/config"
	"hato/mqtt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session           *discordgo.Session
	TriggerEmojiAlias string
	emojiSet          EmojiSet
	publisher         *mqtt.Publisher
}

func MakeBot(conf config.DiscordConfig, publisher *mqtt.Publisher) (*Bot, error) {
	session, err := discordgo.New("Bot " + conf.BotID)
	if err != nil {
		return nil, err
	}

	emojiSet, err := GetEmojiSet()
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:           session,
		TriggerEmojiAlias: conf.TriggerEmojiAlias,
		emojiSet:          *emojiSet,
		publisher:         publisher,
	}, nil
}

func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		return err
	}

	b.session.AddHandler(b.onEmojiAdd)

	return nil
}

func (b *Bot) Stop() error {
	err := b.session.Close()
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) onEmojiAdd(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
	message, err := session.ChannelMessage(event.ChannelID, event.MessageID)
	if err != nil {
		return
	}

	// mentionされたときのみ処理を通す
	me, err := session.User("@me")
	if err != nil {
		return
	}
	if len(message.Mentions) == 0 {
		return
	}
	for i, user := range message.Mentions {
		if user.ID == me.ID {
			break
		}
		if i+1 == len(message.Mentions) {
			return
		}
	}

	// 特定のemojiのときのみ通す
	triggerEmoji, _ := strconv.Unquote(`"` + b.emojiSet.EmojiMap[b.TriggerEmojiAlias] + `"`)
	if event.Emoji.Name != triggerEmoji {
		return
	}

	// メッセージを取得
	str := regexp.MustCompile(`<@\!\d*>`).ReplaceAllString(message.Content, "")
	messageContent := strings.TrimSpace(str)

	for _, reaction := range message.Reactions {
		key := strings.Trim(fmt.Sprintf("%+q", reaction.Emoji.Name), `"`)
		subTopic := b.emojiSet.ReverseEmojiMap[key]
		if subTopic == "" || subTopic == b.TriggerEmojiAlias {
			continue
		}
		b.publisher.Publish(subTopic, messageContent)
	}
}
