package file

import (
  "os"
)

/**
 * Get a channel holding all data written
 * to the given file
 */
func NewReadChannel(file *os.File) (c chan []byte) {
  channel := make(chan []byte)
  go func() {
    data := make([]byte, 1024)
    for {
      bytesRead, err := file.Read(data)
      if err != nil {
        panic(err)
      }
      channel <- data[0:bytesRead]
    }
  }()

  return channel
}
