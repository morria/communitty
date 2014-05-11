package termios

// #include <termios.h>
// #include <unistd.h>
import "C"

func AdjustTermios() {
  struct C.termios stdinTermios
  C.tcgetattr(0, &stdinTermios);
}

