package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dasciam/bedrockscanner/limit"
	"github.com/dasciam/bedrockscanner/output"
	"github.com/dasciam/bedrockscanner/ranges"
	"github.com/dasciam/bedrockscanner/scanner"
	"github.com/samber/lo"
	"log"
	"math"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	var (
		scanWhat         string
		writeToFile      string
		packetsPerSecond int
		numSockets       int
	)

	flag.StringVar(&scanWhat, "what", "ALL", "What to scan (subnet, file path or ALL)")
	flag.IntVar(&packetsPerSecond, "packets-per-second", 5_000, "Number of max packers per second")
	flag.StringVar(&writeToFile, "write-to-file", "", "Path to the file to write results to")
	flag.IntVar(&numSockets, "num-sockets", 1, "Number of sockets")
	flag.Parse()

	log.Printf("Settings:\n- Subnet/file: %s\n- PPS (Packets/s): %d\n- Write to file: %s\n", scanWhat, packetsPerSecond, func() string {
		if writeToFile != "" {
			return writeToFile
		}
		return "(none)"
	}())

	var rng []ranges.Addr

	parsePrefix, err := netip.ParsePrefix(scanWhat)
	if err == nil {
		rng = append(rng, ranges.NewNetIP(parsePrefix))
	} else {
		if strings.ToLower(scanWhat) == "all" {
			const partCount = 64
			const part = math.MaxUint32 / partCount

			for i := 0; i < partCount; i++ {
				if i == partCount-1 {
					rng = append(rng, ranges.NewUInt32(uint32(part*i), math.MaxUint32))
					break
				}
				rng = append(rng, ranges.NewUInt32(uint32(part*i), uint32(part*(i+1))))
			}
			goto next
		}
		data, err := os.ReadFile(scanWhat)
		if err != nil {
			log.Fatal(err)
		}
		rng = lo.FilterMap(bytes.Split(data, []byte("\n")), func(v []byte, _ int) (ranges.Addr, bool) {
			prefix, err := netip.ParsePrefix(string(v))
			if err != nil {
				log.Printf("Error parsing prefix from file %s: %v", scanWhat, err)
				return nil, false
			}
			return ranges.NewNetIP(prefix), true
		})
	}
next:

	var sockets []net.PacketConn

	if numSockets <= 0 {
		panic("got negative number of sockets")
	}
	for range numSockets {
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			panic(err)
		}
		sockets = append(sockets, conn)
	}

	var outputs []scanner.Output

	if writeToFile != "" {
		outputs = append(outputs, output.NewDatabase(writeToFile))
	}
	outputs = append(outputs, output.Print{})

	var (
		wg           sync.WaitGroup
		readWorkerWg sync.WaitGroup
	)

	done := make(chan struct{})

	for _, socket := range sockets {
		readWorkerWg.Add(1)
		go scanner.ReadWorker(socket, output.NewMulti(outputs...), &readWorkerWg, done)
	}

	limiter := limit.NewBasicLimiter(packetsPerSecond)

	for i, r := range rng {
		wg.Add(1)
		go func() {
			scan := scanner.New(r)
			_ = scan.Scan(&wg, sockets[i%len(sockets)], limiter)
		}()
	}
	wg.Wait()
	close(done)
	readWorkerWg.Wait()
	fmt.Println("Ended work")
}
