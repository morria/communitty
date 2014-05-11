package term

/*
#include <termios.h>
#include <unistd.h>
#include <stdio.h>
*/
import "C"

import (
  "errors"
)

// Allows to manipulate the terminal.
type termios struct {
	fd   int // File descriptor
	wrap *_Ctype_struct_termios
}

func NewTermios(fd uintptr) *termios {
  ios := termios{int(fd), new(_Ctype_struct_termios)}
  ios.tcgetattr()
  return &ios;
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
func (tc *termios)KeyPress(fd uintptr) {
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

func (tc *termios)Print() {
  println(tc.wrap.c_iflag)
  println(tc.wrap.c_oflag)
  println(tc.wrap.c_cflag)
  println(tc.wrap.c_lflag)
}

func (tc *termios)Magic() {
  tc.wrap.c_iflag = 11520;
  tc.wrap.c_oflag = 5;
  tc.wrap.c_cflag = 191;
  tc.wrap.c_lflag = 51771;
}

func (tc *termios)Flush() (err error) {
	return tc.tcsetattr(C.TCSAFLUSH);
}

func (tc *termios)MakeRaw() error {
  C.cfmakeraw(tc.wrap);
  return nil
}
