package device

func (m *Manager) TryAcquireSession(sessiongId string) bool {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()

	if m.activeSessionId == "" {
		m.activeSessionId = sessiongId
		return true
	}

	return false
}

func (m *Manager) ReleaseSession(sessionId string) {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()

	if m.activeSessionId == sessionId {
		m.activeSessionId = ""
	}
}

func (m *Manager) CheckSessionAccess(sessionId string) bool {
	m.sessionMu.RLock()
	defer m.sessionMu.RUnlock()

	return m.activeSessionId == sessionId
}
