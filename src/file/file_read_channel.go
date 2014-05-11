package file

import (
  "os"
)

/**
 * Get a channel holding all data written
 * to the given file
 */
func NewReadChannel(file *os.File) (c chan []byte) {
  // Create the channel we'll be writing to
  channel := make(chan []byte)

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

      // Write it to the channel
      channel <- data[0:bytesRead]
    }
  }()

  // Give the channel back while we populate it
  // from the thread
  return channel
}
