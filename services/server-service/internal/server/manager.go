package server

import (
	"io"
	"os/exec"
	"sync"
)

type Manager struct {
	mu      sync.RWMutex
	running map[string]*ServerProcess
}

type ServerStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ServerProcess struct {
	Cmd   *exec.Cmd
	Stdin io.WriteCloser
}

func NewManager() *Manager {
	return &Manager{
		running: make(map[string]*ServerProcess),
	}
}

func (m *Manager) Add(id string, proc *ServerProcess) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running[id] = proc
}

func (m *Manager) Get(id string) (*ServerProcess, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	proc, ok := m.running[id]
	return proc, ok
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.running, id)
}
