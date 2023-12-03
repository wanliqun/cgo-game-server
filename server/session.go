package server

import (
	"net"
	"sync"
	"time"
)

const (
	timeoutCheckInterval = time.Second
	timeOutDuration      = time.Second * 30
)

type Session struct {
	ID         string    // Session ID
	Conn       net.Conn  // Underlying network connection
	LastActive time.Time // Last active time
}

func (s *Session) Close() {
	if s.Conn != nil {
		s.Conn.Close()
	}
}

type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]*Session
	stopChan chan struct{}
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		stopChan: make(chan struct{}),
		sessions: make(map[string]*Session),
	}
}

func (m *SessionManager) Add(sess *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[sess.ID] = sess
}

func (m *SessionManager) Remove(sess *Session) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[sess.ID]; ok {
		delete(m.sessions, sess.ID)
		return true
	}

	return false
}

func (m *SessionManager) CloseAll() {
	for _, s := range m.all() {
		s.Close()
	}
}

func (m *SessionManager) all() (res []*Session) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.sessions {
		res = append(res, s)
	}

	return res
}

// Start function to start checking for timeouts periodically
func (m *SessionManager) Start() {
	ticker := time.NewTicker(timeoutCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.checkTimeout()
		}
	}
}

// Stop function to stop checking for timeouts
func (m *SessionManager) Stop() {
	close(m.stopChan)
}

func (m *SessionManager) checkTimeout() {
	for _, s := range m.all() {
		// Check if the connection has been inactive for longer than the timeout duration.
		if time.Since(s.LastActive) < timeOutDuration {
			continue
		}

		if m.Remove(s) {
			s.Close()

			// TODO: publish session terminated event
		}
	}
}
