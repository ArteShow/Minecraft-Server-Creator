package server

import "sync"

type Manager struct {
	Mu         sync.Mutex
	Containers map[string]string 
}

func NewManager() *Manager {
	return &Manager{
		Containers: make(map[string]string),
	}
}

func (m *Manager) Add(serverID, containerID string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Containers[serverID] = containerID
}

func (m *Manager) Get(serverID string) (string, bool) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	id, ok := m.Containers[serverID]
	return id, ok
}

func (m *Manager) Remove(serverID string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	delete(m.Containers, serverID)
}
