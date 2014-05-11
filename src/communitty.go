package main

// sudo yum install golang golang-github-kr-pty-devel

import (
  "github.com/kr/pty"
  "log"
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

  // Serve HTTP
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

  // Forward input from STDIN to pty master
  channelInput := make(chan []byte)
  go func() {
    data := make([]byte, 1024)
    for {
      bytesRead, err := os.Stdin.Read(data)
      if err != nil {
        panic(err)
      }
      channelInput <- data[0:bytesRead]
      // pty.Write(data);
    }
  }()

  // Forward pty master to STDOUT
  channelOutput := make(chan []byte)
  go func() {
    data := make([]byte, 1024)
    for {
      bytesRead, err := pty.Read(data)
      if err != nil {
        panic(err)
      }
      channelOutput <- data[0:bytesRead]
      // os.Stdout.Write(data);
    }
  }()

  go func() {
    for {
      select {
      case input := <-channelInput:
        pty.Write(input)
        // server.Write(input)
      case output := <-channelOutput:
        os.Stdout.Write(output)
        server.Write(output)
      }
    }
  }()


  err = shell.Wait()
  if nil != err {
    panic(err)
  }
  log.Printf("done");

  /*
  // Listen on websocket at /term
  server := service.NewServer("/term")
  go server.Listen()
  */

  // Tee writes to Stdin out to the server
  // os.Stdin = tee.CreateTee(os.Stdin, server)

  /*
  tty, err := os.OpenFile(os.Stdout.Name(), os.O_RDWR, 600)
  if err != nil {
    panic(err)
  }

  cmd.Stdin = tty
  cmd.Stdout = tty
  cmd.Stderr = tty
  */

  /*
	// cmd.Stdin = io.TeeReader(os.Stdin, server)
  cmd.Stdin = os.Stdin
	cmd.Stdout = io.MultiWriter(os.Stdout, server)
	cmd.Stderr = io.MultiWriter(os.Stderr, server)
  */

  // tty.Write = func (b []byte) (n int, err error) { return f.realFile.Write(b) }


  // tee.TeeTo(os.Stdin, server)

  /*
  go func() {
    stdin, err := cmd.StdoutPipe()
    if err != nil {
      panic(err)
    }

    reader := bufio.NewReader(stdin)
    buf := make([]byte, 128)
    for {
      _, err := reader.Read(buf)
      if err != nil {
        panic(err)
      }
      server.Write(buf)
    }
  }()
  */

  /*
  // Serve static webapp
  http.Handle("/", http.FileServer(http.Dir("./webapp")))

  // Server HTTP
  go func() {
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
      panic("ListenAndServe: " + err.Error())
    }
  }()

  fmt.Println("Serving at http://localhost:3000");
  cmd.Run()
  fmt.Println("No long serving");
  */
}
