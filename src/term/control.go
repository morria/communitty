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

package term

/* Terminal Information and Control. */

// #include <termios.h>
// #include <unistd.h>
import "C"

import "errors"


// Allows to manipulate the terminal.
type termios struct {
	fd   int // File descriptor
	wrap *_Ctype_struct_termios
  isRaw bool
}

func NewTermios(fd int) *termios {
	return &termios{fd, new(_Ctype_struct_termios), false}
}

// Deep copy for pointer fields.
func (tc *termios) CopyTo(to *termios) {
	*to.wrap = *tc.wrap
}

// Gets terminal state.
//
// int tcgetattr(int fd, struct termios *termios_p);
func (tc *termios) tcgetattr() error {
	exitCode, errno := C.tcgetattr(C.int(tc.fd), tc.wrap)

	if exitCode == 0 {
		return nil
	}
	return errno
}

// Sets terminal state.
//
// int tcsetattr(int fd, int optional_actions, const struct termios *termios_p);
func (tc *termios) tcsetattr(optional_actions int) error {
	exitCode, errno := C.tcsetattr(C.int(tc.fd), C.int(optional_actions), tc.wrap)

	if exitCode == 0 {
		return nil
	}
	return errno
}

// === Wrappers around functions.

// Determines if the device is a terminal.
//
// int isatty(int fd);
func Isatty(fd int) bool {
	exitCode, _ := C.isatty(C.int(fd))

	if exitCode == 1 {
		return true
	}
	return false
}

// Determines if the device is a terminal. Return an error, if any.
//
// int isatty(int fd);
func CheckIsatty(fd int) error {
	exitCode, errno := C.isatty(C.int(fd))

	if exitCode == 1 {
		return nil
	}
	return errors.New("it is not a tty: " + errno.Error())
}

// Gets the name of a terminal.
//
// char *ttyname(int fd);
func TTYname(fd int) (string, error) {
	name, errno := C.ttyname(C.int(fd))
	if errno != nil {
		return "", errno
	}
	return C.GoString(name), nil
}

// === Utility

// Turns the echo mode.
func (tc *termios) Echo(echo bool) {
	if !echo {
		tc.wrap.c_lflag &^= (C.ECHO)
	} else {
		tc.wrap.c_lflag |= (C.ECHO)
	}

	if err := tc.tcsetattr(C.TCSAFLUSH); err != nil {
		panic(err)
	}
}

// Sets the terminal to single-character mode.
func (tc *termios)KeyPress(fd int) {
	newSettings := NewTermios(fd)
	tc.CopyTo(newSettings)

	// Disable canonical mode, and set buffer size to 1 byte.
	newSettings.wrap.c_lflag &^= (C.ICANON)
	newSettings.wrap.c_cc[C.VTIME] = 0
	newSettings.wrap.c_cc[C.VMIN] = 1

	if err := newSettings.tcsetattr(C.TCSAFLUSH); err != nil {
		panic("single-character mode")
	}
}

// Sets the terminal to something like the "raw" mode. Input is available
// character by character, echoing is disabled, and all special processing of
// terminal input and output characters is disabled.
//
// Based in C call: void cfmakeraw(struct termios *termios_p)
//
// NOTE: in tty 'raw mode', CR+LF is used for output and CR is used for input.
func (tc *termios)MakeRaw(fd int) error {
	if tc.isRaw {
		return nil
	}

	raw := NewTermios(fd)
	tc.CopyTo(raw)
	tc.isRaw = true

	// Input modes - no break, no CR to NL, no NL to CR, no carriage return,
	// no strip char, no start/stop output control, no parity check.
	raw.wrap.c_iflag &^= (C.BRKINT | C.IGNBRK | C.ICRNL | C.INLCR | C.IGNCR |
		C.ISTRIP | C.IXON | C.PARMRK)

	// Output modes - disable post processing.
	raw.wrap.c_oflag &^= (C.OPOST)

	// Local modes - echoing off, canonical off, no extended functions,
	// no signal chars (^Z,^C).
	raw.wrap.c_lflag &^= (C.ECHO | C.ECHONL | C.ICANON | C.IEXTEN | C.ISIG)

	// Control modes - set 8 bit chars.
	raw.wrap.c_cflag &^= (C.CSIZE | C.PARENB)
	raw.wrap.c_cflag |= (C.CS8)

	// Control chars - set return condition: min number of bytes and timer.
	// We want read to return every single byte, without timeout.
	raw.wrap.c_cc[C.VMIN] = 1 // Read returns when one char is available.
	raw.wrap.c_cc[C.VTIME] = 0

	// Put terminal in raw mode after flushing
	if err := raw.tcsetattr(C.TCSAFLUSH); err != nil {
		return err
	}

	return nil
}

// Restores the original settings for this terminal.
func (tc *termios) RestoreTerm() {
	if tc.isRaw {
		if err := tc.tcsetattr(C.TCSANOW); err != nil {
			panic("restoring the terminal: " + err.Error())
		}
    tc.isRaw = false
	}
}
