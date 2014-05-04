package service

import (
  "code.google.com/p/go.net/websocket"
  "fmt"
)

type Client struct {
  ws *websocket.Conn
  server *Server
  channelMessage chan []byte
}

const channelBufferSize = 256

func NewClient(ws *websocket.Conn, server *Server) *Client {
  if ws == nil {
    panic("ws cannot be nil")
  }

  if server == nil {
    panic("server cannot be nil")
  }

  channelMessage := make(chan []byte, channelBufferSize)

  return &Client{ws, server, channelMessage}
}

func (client *Client) Conn() *websocket.Conn {
  return client.ws
}

func (client *Client) Write(message []byte) {
  select {
  case client.channelMessage <- message:
  default:
    client.server.Remove(client)
    err := fmt.Errorf("client is disconnected")
    if err != nil {
      panic(err)
    }
  }
}

func (client *Client) Listen() {
  for {
    select {
    case message := <-client.channelMessage:
      websocket.Message.Send(client.ws, string(message[:]))
    }
  }
}
