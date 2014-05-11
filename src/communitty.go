package main

import (
  "file"
  "github.com/kr/pty"
  "net/http"
  "os"
  "os/exec"
  "service"
  "term"
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

  go func() {
    channel := term.TrapWinsize()
    select {
    case _ = <-channel:
      row, col := term.GetWinsizeInChar()
      println(row, col)
    }
  }()

  // Run the shell on the pseudo-terminal
  shell := exec.Command(os.Getenv("SHELL"))
  tty, pty, err := pty.Start(shell)
  if err != nil {
    panic(err)
  }

  // Get and set termios properties
  termios := term.Termios(os.Stdin.Fd())
  termios.MakeRaw()
  // termios.DontEcho()
  termios.Flush(os.Stdin.Fd())

  termiosPty := term.Termios(tty.Fd())
  termiosPty.MakeRaw()

  /*
  termios := term.NewTermios(int(pty.Fd()))
  termios.Echo(false)
  termios.MakeRaw(int(pty.Fd()))
  termios.KeyPress(int(pty.Fd()))
  */

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
