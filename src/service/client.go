package service

import (
  "code.google.com/p/go.net/websocket"
  "encoding/base64"
  "encoding/json"
)

/**
 *
 */
type Client struct {
  ws *websocket.Conn
  server *Server
  channelMessage chan []byte
  channelCommand chan string
}

/**
 *
 */
type TerminalMessage struct {
  Command string
  Data string
}

/**
 *
 */
type WindowSizeMessage struct {
  Command string
  Rows int
  Cols int
}

const channelBufferSize = 256

/**
 *
 */
func NewClient(ws *websocket.Conn, server *Server) *Client {
  if ws == nil {
    panic("ws cannot be nil")
  }

  if server == nil {
    panic("server cannot be nil")
  }

  channelMessage := make(chan []byte, channelBufferSize)
  channelCommand := make(chan string)

  return &Client{ws, server, channelMessage, channelCommand}
}

/**
 *
 */
func (client *Client) Conn() *websocket.Conn {
  return client.ws
}

/**
 *
 */
func (client *Client) Write(message []byte) {
  terminalMessage := TerminalMessage{
    "Terminal",
    base64.StdEncoding.EncodeToString(message),
  }

  bytes, err := json.Marshal(terminalMessage)
  if err != nil {
    panic(err)
  }

  client.channelCommand <- string(bytes[:])
}

/**
 * Set the window size
 */
func (client *Client) SetWindowSize(rows, cols uint16) {
  windowSize := WindowSizeMessage{
    "WindowSize",
    int(rows),
    int(cols),
  }

  bytes, err := json.Marshal(windowSize)
  if err != nil {
    panic(err)
  }
  client.channelCommand <- string(bytes[:])
}

/**
 *
 */
func (client *Client) Listen() {
  for {
    select {
    case command := <-client.channelCommand:
      websocket.JSON.Send(client.ws, command)
    }
  }
}
