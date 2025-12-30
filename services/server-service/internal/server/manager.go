package server

type Manager struct {
	running map[string]*ServerProcess
}

func New() *Manager {
	return &Manager{
		running: make(map[string]*ServerProcess),
	}
}

func (m *Manager) Add(id string, proc *ServerProcess) {
	m.running[id] = proc
}

func (m *Manager) Get(id string) (*ServerProcess, bool) {
	proc, ok := m.running[id]
	return proc, ok
}

func (m *Manager) Remove(id string) {
	delete(m.running, id)
}
