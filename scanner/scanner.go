package scanner

import (
	"github.com/dasciam/bedrockscanner/limit"
	"github.com/dasciam/bedrockscanner/raknet"
	"github.com/dasciam/bedrockscanner/ranges"
	"math/rand/v2"
	"net"
	"net/netip"
	"sync"
	"time"
)

type Scanner struct {
	addr ranges.Addr
}

func New(addr ranges.Addr) *Scanner {
	return &Scanner{
		addr: addr,
	}
}

func (s *Scanner) Scan(wg *sync.WaitGroup, conn net.PacketConn, limit limit.Limiter) error {
	var message = s.message()

	lastUpdate := time.Now()

	for {
		addr := s.addr.Addr()

		limit.Increment()
		if time.Since(lastUpdate) > time.Second {
			// Need to update ping time.
			lastUpdate = time.Now()
			message = s.message()
		}
		_, _ = conn.WriteTo(message, net.UDPAddrFromAddrPort(
			netip.AddrPortFrom(addr, 19132),
		))

		var ok bool
		s.addr, ok = s.addr.Next()
		if !ok {
			break
		}
	}
	wg.Done()
	return nil
}

func (s *Scanner) message() []byte {
	return raknet.UnconnectedPing{
		PingTime:   time.Now().UnixNano() / int64(time.Millisecond),
		ClientGUID: rand.Int64(),
	}.Encode()
}
