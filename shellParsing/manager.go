package shellParsing

import (
	"mining-monitoring/log"
)

type Manager struct {
	shellParse *ShellParse
}

func (m *Manager) Run(cmd chan CmdData) {

	go m.shellParse.Receiver(cmd)
	go m.shellParse.Send()
}

func (m *Manager) Close() {
	m.shellParse.Close()
}

func NewManager() (*Manager, error) {
	_, err := log.MyLogicLogger("./log")
	if err != nil {
		return nil, err
	}
	return &Manager{
		shellParse: NewShellParse(),
	}, nil
}
