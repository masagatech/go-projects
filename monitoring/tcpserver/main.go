package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Server ...
type Server struct {
	host      string
	port      string
	heartBeat bool
	clist     *ClientList
}

// Client ...
type Client struct {
	id   string
	conn net.Conn
}

// Config ...
type Config struct {
	Host      string
	Port      string
	HeartBeat bool
}

type ClientList struct {
	mu      sync.Mutex
	clients map[string]*Client
}

// New ...
func New(config *Config) *Server {
	return &Server{
		host:      config.Host,
		port:      config.Port,
		heartBeat: config.HeartBeat,
	}
}

// Run ...
func (server *Server) Run() {

	server.clist = &ClientList{mu: sync.Mutex{}, clients: make(map[string]*Client)}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	// sending heart beat
	if server.heartBeat {
		go server.sendHeart()
	}
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest(server)
	}
}

func (server *Server) sendHeart() {
	for {

		time.Sleep(time.Second * 1)
		for _, v := range server.clist.clients {
			v.conn.Write([]byte("Heart Beat.\n"))
		}

	}
}

func (client *Client) handleRequest(server *Server) {
	server.clist.AddClient(&Client{id: client.conn.LocalAddr().String(), conn: client.conn})
	defer server.clist.RemoveClient(client)

	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {

			fmt.Println("connection closed")
			client.conn.Close()

			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		client.conn.Write([]byte("Message received.\n"))
	}

}

//
func (cList *ClientList) AddClient(cl *Client) {
	cList.mu.Lock()
	defer cList.mu.Unlock()

	cList.clients[cl.conn.RemoteAddr().String()] = cl
	fmt.Println(cList.Size(), cList.clients)
}

func (cList *ClientList) RemoveClient(cl *Client) {
	cList.mu.Lock()
	defer cList.mu.Unlock()

	delete(cList.clients, cl.conn.RemoteAddr().String())

	fmt.Println(cList.Size(), cList.clients)
}

func (cList *ClientList) Size() int {

	return len(cList.clients)
}
