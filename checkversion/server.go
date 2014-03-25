package checkversion

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

type Server struct {
	path           string
	clients        []*Client
	addClient      chan *Client
	removeClient   chan *Client
	versionChanged chan string
}

func NewServer(path string) *Server {
	clients := make([]*Client, 0)
	addClient := make(chan *Client)
	removeClient := make(chan *Client)
	versionChanged := make(chan string)
	return &Server{path, clients, addClient, removeClient, versionChanged}
}

func (self *Server) AddClient() chan<- *Client {
	return (chan<- *Client)(self.addClient)
}

func (self *Server) RemoveClient() chan<- *Client {
	return (chan<- *Client)(self.removeClient)
}

func (self *Server) VersionChanged() chan<- string {
	return (chan<- string)(self.versionChanged)
}

func (self *Server) Listen() {
	log.Println("Start listening...")

	onConnected := func(ws *websocket.Conn) {
		client := NewClient(ws, self)
		self.addClient <- client
		client.Listen()
		defer ws.Close()
	}

	watcher := NewVersionWatcher(self)
	go watcher.Watch("D:\\Applications\\EIP4.0\\Web_Candidates\\modify.notice")

	http.Handle(self.path, websocket.Handler(onConnected))
	log.Println("websocket handler created")

	for {
		select {

		case c := <-self.addClient:
			log.Println("a new client added")
			self.clients = append(self.clients, c)
			log.Println(len(self.clients), "client(s) connected")

		case c := <-self.removeClient:
			log.Println("a new client removed")
			for i := range self.clients {
				if self.clients[i] == c {
					self.clients = append(self.clients[:i], self.clients[i+1:]...)
					log.Println(len(self.clients), "client(s) connected")
					break
				}
			}

		case version := <-self.versionChanged:
			log.Println("version changed, the latest version is:", version, "now notify all clients")
			for _, client := range self.clients {
				client.Notify() <- version
			}
		}
	}
}
