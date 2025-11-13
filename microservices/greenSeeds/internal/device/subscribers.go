package device

func (m *SerialManager) Subscribe() <-chan []byte {
	ch := make(chan []byte, 50)
	m.subsMu.Lock()
	m.subs = append(m.subs, ch)
	m.subsMu.Unlock()

	return ch
}

func (m *SerialManager) Unsubscribe(ch <-chan []byte) {
	m.subsMu.Lock()
	for i, sub := range m.subs {
		if sub == ch {
			m.subs = append(m.subs[:i], m.subs[i+1:]...)
			break
		}
	}
	m.subsMu.Unlock()
}
