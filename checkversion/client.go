package checkversion

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type Client struct {
	ws     *websocket.Conn
	server *Server
	ch     chan string
	done   chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan string)
	done := make(chan bool)

	return &Client{ws, server, ch, done}
}

func (self *Client) Listen() {
	go self.handleWriteRequest()
	self.listenRead()
}

func (self *Client) Notify() chan<- string {
	return (chan<- string)(self.ch)
}

func (self *Client) handleWriteRequest() {
	log.Println("listening write to client")

	for {
		select {

		case version := <-self.ch:
			log.Println("Send:", version)
			websocket.JSON.Send(self.ws, version)

		case <-self.done:
			self.server.RemoveClient() <- self
			return
		}
	}
}

func (self *Client) listenRead() {
	log.Println("listening read from client")

	for {
		select {
		case <-self.done:
			self.server.RemoveClient() <- self
			self.done <- true
			return

		default:
			var version string
			err := websocket.JSON.Receive(self.ws, version)
			if err != nil {
				if err.Error() == "not implemented" {
					// some browsers keep sending 'pong' frame to the server
					// and the websocket library doesn't support it yet
					// so here we simply ignore it
					log.Println("pong frame received, just ignore it here...")
				} else {
					self.done <- true
					return
				}
			} else {
				self.server.VersionChanged() <- version
			}
		}
	}
}
