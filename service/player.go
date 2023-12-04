package service

import (
	"sync"

	"github.com/badu/bus"
	"github.com/wanliqun/cgo-game-server/server"
)

type Player struct {
	Username string
	Session  *server.Session
}

type PlayerService struct {
	mu          sync.Mutex
	usrPlayers  map[string]*Player // username=>Player
	sessPlayers map[string]*Player // session=>Player
}

func NewPlayerService() *PlayerService {
	ps := &PlayerService{
		usrPlayers:  make(map[string]*Player),
		sessPlayers: make(map[string]*Player),
	}
	bus.Sub(ps.OnSessionTerminatedEvent)

	return ps
}

func (s *PlayerService) Add(p *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.usrPlayers[p.Username] = p
	s.sessPlayers[p.Session.ID] = p
}

func (s *PlayerService) Remove(p *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.usrPlayers, p.Username)
	delete(s.sessPlayers, p.Session.ID)
}

func (s *PlayerService) GetByUser(username string) *Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.usrPlayers[username]
}

func (s *PlayerService) GetBySession(sessionID string) *Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.sessPlayers[sessionID]
}

func (s *PlayerService) OnSessionTerminatedEvent(e *server.SessionTerminatedEvent) {
	if player := s.GetBySession(e.Sess.ID); player != nil {
		s.Remove(player)
	}
}
