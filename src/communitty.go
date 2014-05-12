package main

import (
  "github.com/kr/pty"
  "net/http"
  "os"
  "os/exec"
  "service"
  "strconv"
  "term"
)

// Get the original termios so we can reset
// once we're done
var originalTermios = term.GetTermios(os.Stdin.Fd())

/**
 * Reset the terminal, then panic
 */
func panicReset(err error) {
  originalTermios.Flush()
  panic(err)
}

/**
 * Containuously read the given file passing
 * any data read to the given function.
 */
func tail(file *os.File, withData func([]byte)) {
    data := make([]byte, 1024)
    for {
      // Read from the file
      bytesRead, err := file.Read(data)

      // If this thing shuts down, just stop
      // forwarding
      if err != nil {
        return
      }

      withData(data[0:bytesRead])
    }
}

/**
 * Call the given function whenever data comes in
 * on the given channel
 */
func onChannelData(channel chan os.Signal, onData func(os.Signal)) {
  for {
    select {
    case data := <-channel:
      onData(data)
    }
  }
}


/**
 * Make the given file a raw device
 */
func makeRaw(file *os.File) {
  termios := term.GetTermios(file.Fd())
  termios.MakeRaw();
  termios.Flush();
}

/**
 * Serve
 */
func serve(port int, rows, cols uint16) (server *service.Server) {
  // Listen on websocket at /term
  server = service.NewServer("/term", rows, cols)
  go server.Listen()

  // Serve static webapp
  http.Handle("/", http.FileServer(http.Dir("./webapp")))

  go func() {
    err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
    if err != nil {
      panicReset(err)
    }
  }()

  return server
}

/**
 *
 */
func main() {
  // Get the current screen dimensions
  rows, cols := term.GetWindowSize(os.Stdin.Fd())

  // Start serving on port 9000 and listening for clients
  server := serve(9000, rows, cols)

  // Make STDIN a raw device
  makeRaw(os.Stdin)

  // Run the shell on the pseudo-terminal
  shell := exec.Command(os.Getenv("SHELL"))

  // Get a pseudo-terminal and run the command on it
  pty, err := pty.Start(shell)
  if err != nil {
    panicReset(err)
  }

  // Set the initial window size
  term.SetWindowSize(pty.Fd(), rows, cols)

  // Forward window-size changes to the PTY and
  // clients
  go onChannelData(term.TrapWinsize(), func(signal os.Signal) {
    // Get the new size of the window
    newRows, newCols := term.GetWindowSize(os.Stdin.Fd())

    // Set the pseudo-terminal size
    term.SetWindowSize(pty.Fd(), newRows, newCols)

    // Set the size for all clients
    server.SetWindowSize(newRows, newCols)
  })

  // Pipe STDIN to Master
  go tail(os.Stdin, func(data []byte) {
    pty.Write(data)
  })

  // Pipe Master to STDOUT
  go tail(pty, func(data []byte) {
    os.Stdout.Write(data)
    os.Stdout.Sync()
    server.Write(data)
  })

  // Wait for the shell to finish
  err = shell.Wait()
  if nil != err {
    panicReset(err)
  }

  // Reset the original terminal state
  originalTermios.Flush();
}
