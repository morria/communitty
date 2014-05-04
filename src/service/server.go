package service

import (
  "code.google.com/p/go.net/websocket"
  "net/http"
)

/**
 *
 */
type Server struct {
  path          string
  clients       []*Client
  channelAdd    chan *Client
  channelRemove chan *Client
}

/**
 *
 */
func NewServer(path string) *Server {
  clients := make([]*Client, 0)
  channelAdd := make(chan *Client)
  channelRemove := make(chan *Client)

  return &Server{
    path,
    clients,
    channelAdd,
    channelRemove,
  }
}

/**
 *
 */
func (server *Server) Add(client *Client) {
  server.channelAdd <- client
}

/**
 *
 */
func (server *Server) Remove(client *Client) {
  server.channelRemove <- client
}

/**
 *
 */
func (server *Server) Write(message []byte)(int, error) {
  for _, client := range server.clients {
    client.Write(message)
   }
   return len(message), nil
}

/**
 *
 */
func (server *Server) Listen() {

  onConnected := func(ws *websocket.Conn) {

    defer func() {
      err := ws.Close()
      if err != nil {
        panic(err)
      }
    }()

    client := NewClient(ws, server)
    server.Add(client)
    client.Listen()
  }

  http.Handle(server.path, websocket.Handler(onConnected))

  for {
    select {
    case client := <-server.channelAdd:
      server.clients = append(server.clients, client)
    case client := <-server.channelRemove:
      panic(client)
      // server.clients.PushBack(client)
    }
  }

}
