package mux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor interface {
	Process(u *Update) bool
}

type Mux struct {
	Processors []Processor
}

func NewMux() *Mux {
	return &Mux{}
}

func (m *Mux) AddHandler(handlers ...*Handler) *Mux {
	for _, handler := range handlers {
		m.Processors = append(m.Processors, handler)
	}
	return m
}

func (m *Mux) Dispatch(bot *tgbotapi.BotAPI, u tgbotapi.Update) bool {
	return m.Process(&Update{u, bot, false, make(map[string]interface{})})
}

func (m *Mux) Process(u *Update) bool {
	for _, Processor := range m.Processors {
		if Processor.Process(u) {
			return true
		}
	}
	return false
}
