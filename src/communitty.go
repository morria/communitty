package main

import (
  "file"
  "github.com/kr/pty"
  "net/http"
  "os"
  "os/exec"
  "service"
  "termios"
)

func main() {

  // Listen on websocket at /term
  server := service.NewServer("/term")
  go server.Listen()

  // Serve static webapp
  http.Handle("/", http.FileServer(http.Dir("./webapp")))

  // Serve HTTP at http://localhost:3000/
  go func() {
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
      panic("ListenAndServe: " + err.Error())
    }
  }()

  // Get and set termios properties
  termios := termios.Termios()
  termios.MakeRaw()
  termios.DontEcho()
  termios.TCSAFlush(os.Stdin.Fd())

  // Run the shell on the pseudo-terminal
  shell := exec.Command(os.Getenv("SHELL"))
  pty, err := pty.Start(shell)
  if err != nil {
    panic(err)
  }

  // Read from input and output channels, writing to
  // the local TTY and any websocket clients
  go func() {
    // Get all data written to STDIN on a channel
    channelInput := file.NewReadChannel(os.Stdin)

    // Get all data written to the PTY on a channel
    channelOutput := file.NewReadChannel(pty)

    // Forward data to the children appropriately
    for {
      select {
      case input := <-channelInput:
        pty.Write(input)
      case output := <-channelOutput:
        os.Stdout.Write(output)
        server.Write(output)
      }
    }
  }()

  // Wait for the shell to finish
  err = shell.Wait()
  if nil != err {
    panic(err)
  }
}
