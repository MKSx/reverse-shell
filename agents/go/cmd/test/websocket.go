package test

import (
	"bytes"
	"net"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/maxlaverse/reverse-shell/message"
	"github.com/maxlaverse/reverse-shell/util"
)

type DummyMaster struct {
	agentConn  convert2
	t          *testing.T
	port       int32
	messageCh  chan interface{}
	messageCh2 chan message.Serializable
}

func NewDummyMaster(t *testing.T) DummyMaster {
	b := DummyMaster{
		t:          t,
		messageCh:  make(chan interface{}),
		messageCh2: make(chan message.Serializable),
	}

	go http.HandleFunc("/agent/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := util.WebSocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Could not upgrade the connection: %v", err)
		}
		go func() {
			defer conn.Close()
			for {
				_, m, err := conn.ReadMessage()
				if err != nil {
					t.Errorf("Could read message: %v", err)
				}
				b.messageCh <- message.FromBinary(m)
			}
		}()

		go func() {
			for {
				m := <-b.messageCh2
				conn.WriteMessage(websocket.BinaryMessage, message.ToBinary(m))
			}

		}()
	})

	// Start a fake master
	tcpMaster, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error start dummy tcp master: %v", err)
	}

	port, err := AddrToPort(tcpMaster.Addr())
	if err != nil {
		t.Fatalf("Error while getting dummy tcp master port from '%s': %v", tcpMaster.Addr(), err)
	}

	b.port = port

	go func() {
		err = http.Serve(tcpMaster, nil)
		if err != nil {
			t.Fatalf("Error listening '%s': %v", tcpMaster.Addr(), err)
		}
	}()

	return b
}

func (b *DummyMaster) ReadRawMessage() interface{} {
	return <-b.messageCh
}

type convert2 func(messageType int, data []byte) error

func (b *DummyMaster) SendMessage(m message.Serializable) {
	b.messageCh2 <- m
}

func (b *DummyMaster) Port() int32 {
	return b.port
}

func (b *DummyMaster) ReadMessageUntilTerminated() string {
	inputBuf := bytes.NewBufferString("")
	for {
		m := <-b.messageCh
		switch v := m.(type) {
		case *message.ProcessOutput:
			inputBuf.Write(v.Data)
		case *message.ProcessTerminated:
			return inputBuf.String()
		default:
			b.t.Errorf("Unexpected message type: %v", b)
		}
	}
}
