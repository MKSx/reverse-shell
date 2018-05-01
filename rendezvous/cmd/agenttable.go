package cmd

import "github.com/gorilla/websocket"

type AgentTable struct {
	agents map[string]*websocket.Conn
}

func NewAgentTable() AgentTable {
	return AgentTable{
		agents: make(map[string]*websocket.Conn),
	}
}

func (s *AgentTable) AddAgent(conn *websocket.Conn) {
	s.agents[conn.RemoteAddr().String()] = conn
}

func (s *AgentTable) RemoveAgent(conn *websocket.Conn) {
	delete(s.agents, conn.RemoteAddr().String())
}

func (b *AgentTable) FindAgent(address string) *websocket.Conn {
	return b.agents[address]
}
