package raknet

import (
	"encoding/binary"
	"io"
)

type UnconnectedPong struct {
	PingTime   int64
	ServerGUID int64
	Data       []byte
}

func (pk *UnconnectedPong) Decode(data []byte) error {
	if len(data) < 34 || len(data) < 34+int(binary.BigEndian.Uint16(data[32:])) {
		return io.ErrUnexpectedEOF
	}
	pk.PingTime = int64(binary.BigEndian.Uint64(data))
	pk.ServerGUID = int64(binary.BigEndian.Uint64(data[8:]))
	// Magic: 16 bytes.
	n := binary.BigEndian.Uint16(data[32:])
	pk.Data = append([]byte(nil), data[34:34+n]...)
	return nil
}
