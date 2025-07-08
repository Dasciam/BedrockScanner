package scanner

import "net"

type Output interface {
	Write(addr net.Addr, pong Pong)
}
