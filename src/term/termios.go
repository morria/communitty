package term

// #include <termios.h>
// #include <unistd.h>
import "C"

/**
 *
 */
type Termios struct {
	fd   int // File descriptor
	wrap *_Ctype_struct_termios
}

/**
 *
 */
func GetTermios(fd uintptr) *Termios {
  ios := Termios{int(fd), new(_Ctype_struct_termios)}
  ios.tcgetattr()
  return &ios;
}

/**
 *
 */
func (tc *Termios) CopyTo(to *Termios) {
	*to.wrap = *tc.wrap
}

/**
 *
 */
func (tc *Termios) tcgetattr() error {
	exitCode, errno := C.tcgetattr(C.int(tc.fd), tc.wrap)

	if exitCode == 0 {
		return nil
	}
	return errno
}

/**
 *
 */
func (tc *Termios) tcsetattr(optional_actions int) error {
	exitCode, errno := C.tcsetattr(C.int(tc.fd), C.int(optional_actions), tc.wrap)

	if exitCode == 0 {
		return nil
	}
	return errno
}

func (tc *Termios)Print() {
  println(tc.wrap.c_iflag)
  println(tc.wrap.c_oflag)
  println(tc.wrap.c_cflag)
  println(tc.wrap.c_lflag)
}

func (tc *Termios)Magic() {
  tc.wrap.c_iflag = 11520;
  tc.wrap.c_oflag = 5;
  tc.wrap.c_cflag = 191;
  tc.wrap.c_lflag = 51771;
}

func (tc *Termios)Flush() (err error) {
	return tc.tcsetattr(C.TCSAFLUSH);
}

func (tc *Termios)MakeRaw() error {
  C.cfmakeraw(tc.wrap);
  return nil
}
