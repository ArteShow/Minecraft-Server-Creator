package server

import (
	"os/exec"
)

type Manager struct {
	running map[string]*exec.Cmd
}

func New() *Manager {
	return &Manager{
		running: make(map[string]*exec.Cmd),
	}
}

func (m *Manager) Add(id string, cmd *exec.Cmd) {
	m.running[id] = cmd
}

func (m *Manager) Get(id string) (*exec.Cmd, bool) {
	cmd, ok := m.running[id]
	return cmd, ok
}

func (m *Manager) Remove(id string) {
	delete(m.running, id)
}
