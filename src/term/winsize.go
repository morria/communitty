// Copyright 2010  The "Go-Term" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 == Window and Terminal Size Services

 The kernel stores the winsize structure to provide a consistent interface for
 the current terminal or window size.

 The winsize structure contains the following fields:

 + ws_row: Indicates the number of rows (in characters) on the window or terminal
 + ws_col: Indicates the number of columns (in characters) on the window or terminal
 + ws_xpixel: Indicates the horizontal size (in pixels) of the window or terminal
 + ws_ypixel: Indicates the vertical size (in pixels) of the window or terminal

 === Queries terminal characteristics

 + TIOCGWINSZ: Gets the window size. The argument to this ioctl operation is a
 pointer to a winsize structure, into which the current terminal or window size
 is placed.

 + TIOCSWINSZ: Sets the window size. The argument to this ioctl operation is a
 pointer to a winsize structure, which is used to set the current terminal or
 window size information. If the new information differs from the previous, a
 SIGWINCH signal is sent to the terminal process group.
*/

package term

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unsafe"
)

// Values by default
const (
	_ROW    = 24
	_COLUMN = 80
)

var ChanWinsize = make(chan int) // Allocate a channel for TrapWinsize()

// === Get

// Gets the window size using the TIOCGWINSZ ioctl() call on the tty device.
func GetWinsize() (*winsize, error) {
	ws := new(winsize)

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}

// Gets the number of rows and columns (in characters) on the window or terminal.
func GetWinsizeInChar() (row, col int) {
	ws, err := GetWinsize()

	// If there is any error, then to try get the values through environment.
	// Else, it is used values by default.
	if err != nil {
		sRow := os.Getenv("LINES")
		sCol := os.Getenv("COLUMNS")

		if sRow == "" {
			row = _ROW
		} else {
			iRow, err := strconv.Atoi(sRow)
			if err == nil {
				row = iRow
			} else {
				row = _ROW
			}
		}

		if sCol == "" {
			col = _COLUMN
		} else {
			iCol, err := strconv.Atoi(sCol)
			if err == nil {
				col = iCol
			} else {
				col = _COLUMN
			}
		}
		return
	}

	return int(ws.Row), int(ws.Col)
}

// Caughts a signal named SIGWINCH whenever the screen size changes.
func TrapWinsize() (chan os.Signal) {
  channel := make(chan os.Signal)
  signal.Notify(channel, syscall.SIGWINCH)
  return channel

  /*
	var sig os.Signal
	go func() {
		for sig = range signal.Incoming {
			// if sig.(os.UnixSignal) == syscall.SIGWINCH {
			if sig == syscall.SIGWINCH {
				ChanWinsize <- 1 // Send a signal
			}
		}
	}()
  */
}
