package service

import (
	"sync"

	"github.com/badu/bus"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/server"
)

const (
	CtxKeyPlayer server.ContextKey = "player"
)

type Player struct {
	Username string
	Session  *server.Session
}

type PlayerService struct {
	mu          sync.Mutex
	usrPlayers  map[string]*Player // username=>Player
	sessPlayers map[string]*Player // session=>Player
	sessionMgr  *server.SessionManager
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

func (s *PlayerService) Kickoff(p *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.usrPlayers, p.Username)
	delete(s.sessPlayers, p.Session.ID)

	s.sessionMgr.Terminate(p.Session)
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

func (s *PlayerService) Login(req *proto.LoginRequest, session *server.Session) (*Player, error) {
	// TODO check password with the configured server password.

	player := s.GetByUser(req.Username)

	if player != nil && player.Session.ID == session.ID {
		// User already logined with the same session.
		return player, nil
	}

	if player != nil {
		// Kick off the player with an old session.
		s.Kickoff(player)
	}

	player = &Player{
		Username: req.Username,
		Session:  session,
	}
	s.Add(player)

	return player, nil
}

func (s *PlayerService) OnSessionTerminatedEvent(e *server.SessionTerminatedEvent) {
	if player := s.GetBySession(e.Sess.ID); player != nil {
		s.Kickoff(player)
	}
}
