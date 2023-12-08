package server

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/badu/bus"
	"github.com/google/uuid"
)

const (
	timeoutCheckInterval = time.Second
	timeOutDuration      = time.Second * 30

	CtxKeySession ContextKey = "session"
)

func NewContextFromSession(parent context.Context, sess *Session) context.Context {
	return context.WithValue(parent, CtxKeySession, sess)
}

func SessionFromContext(ctx context.Context) (sess *Session, ok bool) {
	sess, ok = ctx.Value(CtxKeySession).(*Session)
	return sess, ok
}

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

func (m *SessionManager) Terminate(sess *Session) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.sessions[sess.ID]; ok {
		s.Close()
		delete(m.sessions, sess.ID)
		return true
	}

	return false
}

func (m *SessionManager) TerminateAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.sessions {
		s.Close()
		delete(m.sessions, s.ID)
	}
}

func (m *SessionManager) ListAll() []*Session {
	return m.all()
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
		// Check if the session has been inactive for longer than the timeout duration.
		if time.Since(s.LastActive()) >= timeOutDuration {
			m.Terminate(s)
			// Publish session terminated event
			bus.Pub(&SessionTerminatedEvent{Sess: s})
		}
	}
}
