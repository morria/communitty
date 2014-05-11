package termios

// #include <termios.h>
import "C"

func Termios() (t *_Ctype_struct_termios) {
  ios := new(_Ctype_struct_termios)
  C.tcgetattr(0, ios)
  return ios
}

func (ios *_Ctype_struct_termios) MakeRaw() {
  C.cfmakeraw(ios);
}

func (ios *_Ctype_struct_termios)DontEcho() {
  ((*ios).c_lflag) = (((*ios).c_lflag) & ^C.tcflag_t(C.ECHO))
}

func (ios *_Ctype_struct_termios) TCSAFlush(fd uintptr) {
  C.tcsetattr(C.int(fd), C.TCSAFLUSH, ios);
}
