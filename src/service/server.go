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
  windowRows    uint16
  windowCols    uint16
}

/**
 *
 */
func NewServer(path string, rows, cols uint16) *Server {
  clients := make([]*Client, 0)
  channelAdd := make(chan *Client)
  channelRemove := make(chan *Client)

  return &Server{
    path,
    clients,
    channelAdd,
    channelRemove,
    rows,
    cols,
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
func (server *Server) SetWindowSize(rows, cols uint16)(error) {
  for _, client := range server.clients {
    client.SetWindowSize(rows, cols)
  }
  return nil
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
      client.SetWindowSize(
        server.windowRows, server.windowCols)
      println("client connected");
    case client := <-server.channelRemove:
      panic(client)
      // server.clients.PushBack(client)
    }
  }

}
