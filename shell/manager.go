package shell

type Manager struct {
	shellParse *Parse
}

func (m *Manager) Run(cmd chan CmdData) {

	for i := 0; i < 100; i++ {
		go m.shellParse.Receiver(cmd)
	}
	go m.shellParse.Send()
}

func (m *Manager) Close() {
	m.shellParse.Close()
}

func NewManager() (*Manager, error) {
	return &Manager{
		shellParse: NewShellParse(),
	}, nil
}
