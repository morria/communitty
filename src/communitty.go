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
  // Get the current screen dimensions
  rows, cols := term.GetWindowSize(os.Stdin.Fd())

  // Listen on websocket at /term
  server := service.NewServer("/term", rows, cols)
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


  // Run the shell on the pseudo-terminal
  shell := exec.Command(os.Getenv("SHELL"))

  // Get a pseudo-terminal and run the command
  // on it
  pty, err := pty.Start(shell)
  if err != nil {
    panic(err)
  }

  // Make STDIN a raw device
  termios := term.Termios(os.Stdin.Fd())
  termios.MakeRaw()
  termios.Flush(os.Stdin.Fd())

  // Forward window-size changes to the PTY and
  // clients
  go func() {
    channel := term.TrapWinsize()
    select {
    case _ = <-channel:
      // Get the new size of the window
      newRows, newCols := term.GetWindowSize(os.Stdin.Fd())

      // Set the pseudo-terminal size
      term.SetWindowSize(pty.Fd(), newRows, newCols)

      // Set the size for all clients
      server.SetWindowSize(newRows, newCols)
    }
  }()

  // Set the initial window size
  term.SetWindowSize(pty.Fd(), rows, cols)

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
