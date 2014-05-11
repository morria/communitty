package main

import (
  "github.com/kr/pty"
  "net/http"
  "os"
  "os/exec"
  "service"
  "term"
  "syscall"
)

/**
 * Containuously read the given file passing
 * any data read to the given function.
 */
func tail(file *os.File, withData func([]byte)) {
  go func() {
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
  }()
}

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

  // Get the original termios so we can reset
  // once we're done
  originalTermios := term.GetTermios(os.Stdin.Fd())

  // Run the shell on the pseudo-terminal
  shell := exec.Command(os.Getenv("SHELL"))

  /*
  // Get a pseudo-terminal and run the command
  // on it
  tty, pty, err := pty.Start(shell)
  if err != nil {
    panic(err)
  }
  */

	ptty, tty, err := pty.Open()
	if err != nil {
    panic(err)
	}

  slaveTermios := term.GetTermios(tty.Fd())
  originalTermios.CopyTo(slaveTermios)
  term.SetWindowSize(tty.Fd(), rows, cols)
  slaveTermios.Flush()

	defer tty.Close()
	shell.Stdout = tty
	shell.Stdin = tty
	shell.Stderr = tty

	shell.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	err = shell.Start()
	if err != nil {
    panic(err)
	}


  // Make STDIN a raw terminal
  termios := term.GetTermios(os.Stdin.Fd())
  termios.MakeRaw();
  // termios.Echo(false);
  // termios.Magic();
  termios.Flush();

  // Forward window-size changes to the PTY and
  // clients
  go func() {
    channel := term.TrapWinsize()
    select {
    case _ = <-channel:
      // Get the new size of the window
      newRows, newCols := term.GetWindowSize(os.Stdin.Fd())

      // Set the pseudo-terminal size
      term.SetWindowSize(ptty.Fd(), newRows, newCols)

      // Set the size for all clients
      server.SetWindowSize(newRows, newCols)
    }
  }()

  // Set the initial window size
  term.SetWindowSize(ptty.Fd(), rows, cols)

  tail(os.Stdin, func(data []byte) {
    ptty.Write(data)
  })

  tail(ptty, func(data []byte) {
    os.Stdout.Write(data)
    os.Stdout.Sync()
    server.Write(data)
  })

  // Wait for the shell to finish
  err = shell.Wait()
  if nil != err {
    panic(err)
  }

  // Reset the original terminal state
  originalTermios.Flush();
}
