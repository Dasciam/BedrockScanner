package scanner

import (
	"errors"
	"fmt"
	"github.com/dasciam/bedrockscanner/raknet"
	"net"
	"os"
	"sync"
	"time"
)

func ReadWorker(conn net.PacketConn, output Output, wg *sync.WaitGroup, done <-chan struct{}) {
	var (
		buffer = make([]byte, 1500)
		pong   = new(raknet.UnconnectedPong)
	)

	defer wg.Done()

	go func() {
		<-done
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	}()

	for {
		n, addr, err := conn.ReadFrom(buffer)

		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Printf("Read deadline exceeded!\n")
				return
			}
			fmt.Printf("Error reading from socket: %v\n", err)
			continue
		}
		if n == 0 || buffer[0] != 0x1C {
			// Non-pong packet received. Ignore.
			continue
		}
		err = pong.Decode(buffer[1:n])
		if err != nil {
			fmt.Printf("Error decoding packet: %v\n", err)
			continue
		}
		data, ok := PongFromBytes(pong.Data)
		if !ok {
			continue
		}
		output.Write(addr, data)
	}
}
