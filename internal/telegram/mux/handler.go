package mux

import (
	"regexp"
	"strings"
)

type HandleFunc func(u *Update)

type Handler struct {
	Filter  FilterFunc
	Handles []HandleFunc
}

func (h *Handler) Process(u *Update) bool {
	if h.Filter(u) {
		for i := 0; i < len(h.Handles) && !u.Consumed; i++ {
			h.Handles[i](u)
		}
		return true
	}
	return false
}

func NewHandler(filter FilterFunc, handles []HandleFunc) *Handler {
	if filter == nil {
		filter = Any()
	}
	return &Handler{filter, handles}
}

func NewMessageHandler(filter FilterFunc, handles []HandleFunc) *Handler {
	newFilter := IsMessage()
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewHandler(newFilter, handles)
}

func NewCommandHandler(cmd string, filter FilterFunc, handles ...HandleFunc) *Handler {
	handles = append([]HandleFunc{
		func(u *Update) {
			u.Context["args"] = strings.Split(u.Message.Text, " ")[1:]
		},
	}, handles...)

	commandsFilters := []FilterFunc{}

	for _, variant := range strings.Split(cmd, " ") {
		commandsFilters = append(commandsFilters, IsCommandMessage(variant))
	}
	newFilter := Or(commandsFilters...)
	if filter != nil {
		filter = And(newFilter, filter)
	}
	return NewMessageHandler(filter, handles)
}

func NewCallBackQuertHandler(pattern string, filter FilterFunc, handles ...HandleFunc) *Handler {
	exp := regexp.MustCompile(pattern)
	newFilter := And(IsCallBackQuery(), func(u *Update) bool {
		return exp.Match([]byte(u.CallbackQuery.Data))
	})
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	handles = append([]HandleFunc{
		func(u *Update) {
			u.Context["exp"] = exp
			u.Context["matches"] = exp.FindStringSubmatch(u.CallbackQuery.Data)
		},
	}, handles...)
	return NewHandler(newFilter, handles)
}
