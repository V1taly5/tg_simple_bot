package menu

import "simple_tg_bot/internal/telegram/models"

type Menu struct {
	Commands map[string]models.CmdFunc
}

func New() *Menu {
	cmds := make(map[string]models.CmdFunc)

	return &Menu{
		Commands: cmds,
	}
}

func (m *Menu) RegisterCmd(cmd string, CmdFunc models.CmdFunc) {
	if m.Commands == nil {
		m.Commands = make(map[string]models.CmdFunc)
	}
	m.Commands[cmd] = CmdFunc
}

// Command get a command and ok from menu.
func (m *Menu) Command(cmd string) (models.CmdFunc, bool) {
	cf, ok := m.Commands[cmd]
	return cf, ok
}
