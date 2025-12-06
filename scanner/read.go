package scanner

import (
	"errors"
	"github.com/dasciam/bedrockscanner/raknet"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var id uint32

func ReadWorker(conn net.PacketConn, output Output, wg *sync.WaitGroup, timeout time.Duration) {
	var (
		buffer = make([]byte, 1500)
		pong   = new(raknet.UnconnectedPong)
	)

	wId := atomic.AddUint32(&id, 1)

	defer wg.Done()

	for {
		_ = conn.SetReadDeadline(time.Now().Add(timeout))
		n, addr, err := conn.ReadFrom(buffer)

		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				log.Printf("%d: read deadline exceeded", wId)
				return
			}
			log.Printf("%d: error reading from socket: %v", wId, err)
			continue
		}
		if n == 0 || buffer[0] != 0x1C {
			// Non-pong packet received. Ignore.
			continue
		}
		err = pong.Decode(buffer[1:n])
		if err != nil {
			log.Printf("%d: error decoding packet: %v", wId, err)
			continue
		}
		data, ok := PongFromBytes(pong.Data)
		if !ok {
			continue
		}
		output.Write(addr, data)
	}
}
