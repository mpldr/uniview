package mpv

import (
	"bytes"
	"net"
)

func readUntilLF(conn net.Conn) ([]byte, error) {
	b := make([]byte, 1)
	var result bytes.Buffer

	for {
		_, err := conn.Read(b)
		if err != nil {
			return nil, err
		}

		result.Write(b)

		if b[0] == '\n' {
			break
		}
	}

	return result.Bytes(), nil
}
