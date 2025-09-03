package internal

import (
	"io"
	"net"
)

// closeWrite attempts a half-close on TCP; no-op for non-TCP conns.
func closeWrite(c net.Conn) {
	if tc, ok := c.(*net.TCPConn); ok {
		_ = tc.CloseWrite()
	}
}

func Pipe(a, b net.Conn) {
	defer a.Close()
	defer b.Close()

	errc := make(chan error, 2)

	// a -> b
	go func() {
		_, err := io.Copy(b, a)
		closeWrite(b) // signal EOF downstream
		errc <- err
	}()

	// b -> a
	go func() {
		_, err := io.Copy(a, b)
		closeWrite(a)
		errc <- err
	}()

	// wait for both directions to finish (or first hard error)
	<-errc
	<-errc
}
