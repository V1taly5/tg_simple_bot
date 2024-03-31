package mux

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update struct {
	tgbotapi.Update
	Bot      *tgbotapi.BotAPI
	Consumed bool
	Context  map[string]interface{}
}

func (u *Update) Consume() {
	u.Consumed = true
}

// GetUser retrieves User object from update.
func (u *Update) GetUser() *tgbotapi.User {
	if u.Message != nil {
		return u.Message.From
	} else if u.EditedMessage != nil {
		return u.EditedMessage.From
	} else if u.ChannelPost != nil {
		return u.ChannelPost.From
	} else if u.EditedChannelPost != nil {
		return u.EditedMessage.From
	} else if u.InlineQuery != nil {
		return u.InlineQuery.From
	} else if u.ChosenInlineResult != nil {
		return u.ChosenInlineResult.From
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From
	} else if u.ShippingQuery != nil {
		return u.ShippingQuery.From
	} else if u.PreCheckoutQuery != nil {
		return u.PreCheckoutQuery.From
	}
	log.Println("object was not found in update")
	return nil
}

// GetMessage retrieves message object from update
func (u *Update) GetMesssge() *tgbotapi.Message {
	candidates := []*tgbotapi.Message{u.Message, u.EditedMessage, u.ChannelPost, u.EditedChannelPost}
	for _, message := range candidates {
		if message != nil {
			return message
		}
	}

	if u.CallbackQuery != nil {
		return u.CallbackQuery.Message
	}
	return nil
}

func (u *Update) GetChat() *tgbotapi.Chat {
	message := u.GetMesssge()
	if message != nil {
		return message.Chat
	}
	return nil
}
