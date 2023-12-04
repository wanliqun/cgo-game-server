package server

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	timeoutCheckInterval = time.Second
	timeOutDuration      = time.Second * 30

	CtxKeySession ContextKey = "session"
)

type Session struct {
	ID         string   // Session ID
	Conn       net.Conn // Underlying network connection
	lastActive int64    // Last active timestamp
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		ID:         uuid.NewString(),
		Conn:       conn,
		lastActive: time.Now().Unix(),
	}
}

func (s *Session) Refresh() {
	atomic.StoreInt64(&s.lastActive, time.Now().Unix())
}

func (s *Session) LastActive() time.Time {
	return time.Unix(atomic.LoadInt64(&s.lastActive), 0)
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

func (m *SessionManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions = make(map[string]*Session)
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
		if time.Since(s.LastActive()) < timeOutDuration {
			continue
		}

		if m.Remove(s) {
			s.Close()

			// TODO: publish session terminated event
		}
	}
}
