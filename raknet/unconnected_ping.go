package raknet

import "encoding/binary"

type UnconnectedPing struct {
	PingTime   int64
	ClientGUID int64
}

func (ping UnconnectedPing) Encode() (data []byte) {
	b := make([]byte, 33)
	b[0] = 1
	binary.BigEndian.PutUint64(b[1:], uint64(ping.PingTime))
	copy(b[9:], magic[:])
	binary.BigEndian.PutUint64(b[25:], uint64(ping.ClientGUID))
	return b
}
