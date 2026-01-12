package server

import "sync"

type Manager struct {
	mu         sync.Mutex
	containers map[string]string 
}

func NewManager() *Manager {
	return &Manager{
		containers: make(map[string]string),
	}
}

func (m *Manager) Add(serverID, containerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containers[serverID] = containerID
}

func (m *Manager) Get(serverID string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.containers[serverID]
	return id, ok
}

func (m *Manager) Remove(serverID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.containers, serverID)
}
