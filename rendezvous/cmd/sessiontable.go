package cmd

import "github.com/gorilla/websocket"

type SessionState int

const (
	SESSION_OPEN SessionState = 1 + iota
	SESSION_CLOSED
	SESSION_LOST
)

type Session struct {
	Id         string
	masterConn []*websocket.Conn
	agentConn  *websocket.Conn
	State      SessionState
}

type SessionTable struct {
	sessionTable map[string]*Session
}

func (s SessionState) String() string {
	switch s {
	case SESSION_LOST:
		return "lost"
	case SESSION_CLOSED:
		return "closed"
	case SESSION_OPEN:
		return "open"
	}
	return "unknown"
}

func NewSessionTable() SessionTable {
	return SessionTable{
		sessionTable: make(map[string]*Session),
	}
}

func (s *SessionTable) AddSession(sess *Session) {
	s.sessionTable[sess.Id] = sess
}

func (s *SessionTable) AttachToSession(id string, conn *websocket.Conn) {
	s.sessionTable[id].masterConn = append(s.sessionTable[id].masterConn, conn)
}

func (s *SessionTable) FindSession(session string) *Session {
	return s.sessionTable[session]
}

func (s *SessionTable) FindSessionByAgent(sconn *websocket.Conn) []*Session {
	var w []*Session
	for _, v := range s.sessionTable {
		if v.agentConn == sconn {
			w = append(w, v)
		}
	}
	return w
}
