package term

// #include <sys/ioctl.h>
// #include <termios.h>
import "C"

/**
 * Get termios for the given file descriptor
 */
func Termios(fd uintptr) (t *_Ctype_struct_termios) {
  ios := new(_Ctype_struct_termios)
  C.tcgetattr(C.int(fd), ios)
  return ios
}

/**
 * Make the given termios "raw"
 */
func (ios *_Ctype_struct_termios) MakeRaw() {
  C.cfmakeraw(ios);
}

/**
 *
 */
func (ios *_Ctype_struct_termios)DontEcho() {
  ((*ios).c_lflag) = (((*ios).c_lflag) & ^C.tcflag_t(C.ECHO))
}

/**
 *
 */
func (ios *_Ctype_struct_termios) Flush(fd uintptr) {
  C.tcsetattr(C.int(fd), C.TCSAFLUSH, ios);
}
