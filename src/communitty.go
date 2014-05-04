package main

// sudo yum install golang golang-github-kr-pty-devel

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "os/exec"
  "service"
  // "github.com/kr/pty"
  // "log"
  // "syscall"
)

func main() {
  cmd := exec.Command("/bin/zsh")

  /*
	_, tty, err := pty.Open()
	if err != nil {
    panic(err);
	}
	defer tty.Close()
  */

  // Listen on websocket at /term
  server := service.NewServer("/term")
  go server.Listen()

	cmd.Stdin = os.Stdin
	cmd.Stdout = io.MultiWriter(os.Stdout, server)
	cmd.Stderr = io.MultiWriter(os.Stderr, server)

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
}
