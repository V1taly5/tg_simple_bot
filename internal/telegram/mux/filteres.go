package mux

import "regexp"

// FilterFunc is employed to determine whether this update should be handled by the handler.
type FilterFunc func(u *Update) bool

var commandRegex = regexp.MustCompile("^/([0-9a-zA-Z_]+)(@[0-9a-zA-Z_]{3,})?")

func Any() FilterFunc {
	return func(u *Update) bool {
		return true
	}
}

func IsMessage() FilterFunc {
	return func(u *Update) bool {
		return u.Message != nil
	}
}

func IsInlineQuery() FilterFunc {
	return func(u *Update) bool {
		return u.InlineQuery != nil
	}
}

func IsCallBackQuery() FilterFunc {
	return func(u *Update) bool {
		return u.CallbackQuery != nil
	}
}

func IsEditedMessage() FilterFunc {
	return func(u *Update) bool {
		return u.EditedMessage != nil
	}
}

func IsChannelPost() FilterFunc {
	return func(u *Update) bool {
		return u.ChannelPost != nil
	}
}

func IsEditedChannelPost() FilterFunc {
	return func(u *Update) bool {
		return u.EditedChannelPost != nil
	}
}

func HasText() FilterFunc {
	return func(u *Update) bool {
		message := u.GetMesssge()
		return message != nil && message.Text != "" && message.Text[0] != '/'
	}
}

func IsAnyCommandMessage() FilterFunc {
	return And(IsMessage(), func(u *Update) bool {
		matches := commandRegex.FindStringSubmatch(u.Message.Text)
		if len(matches) == 0 {
			return false
		}
		botName := matches[2]
		if botName != "" && botName != "@"+u.Bot.Self.UserName {
			return false
		}
		return true
	})
}

func IsCommandMessage(cmd string) FilterFunc {
	return And(IsAnyCommandMessage(), func(u *Update) bool {
		matches := commandRegex.FindStringSubmatch(u.Message.Text)
		actualCmd := matches[1]
		return actualCmd == cmd
	})
}

func HasPhoto() FilterFunc {
	return func(u *Update) bool {
		msg := u.GetMesssge()
		return msg != nil && msg.Photo != nil
	}
}

func HasDocument() FilterFunc {
	return func(u *Update) bool {
		msg := u.GetMesssge()
		return msg != nil && msg.Document != nil
	}
}

func And(filters ...FilterFunc) FilterFunc {
	return func(u *Update) bool {
		for _, filter := range filters {
			if !filter(u) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...FilterFunc) FilterFunc {
	return func(u *Update) bool {
		for _, filter := range filters {
			if filter(u) {
				return true
			}
		}
		return false
	}
}

func Not(filter FilterFunc) FilterFunc {
	return func(u *Update) bool {
		return !filter(u)
	}
}
