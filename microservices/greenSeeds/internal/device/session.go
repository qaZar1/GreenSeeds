package device

func (m *Manager) TryAcquireSession(sessiongId string) bool {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()

	if m.activeSessionId == "" {
		m.activeSessionId = sessiongId
		m.log.Info().Msg("session acquired")
		return true
	}

	return false
}

func (m *Manager) ReleaseSession(sessionId string) {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()

	if m.activeSessionId == sessionId {
		m.activeSessionId = ""
		m.log.Info().Msg("session released")
	}
}

func (m *Manager) CheckSessionAccess(sessionId string) bool {
	m.sessionMu.RLock()
	defer m.sessionMu.RUnlock()

	return m.activeSessionId == sessionId
}
