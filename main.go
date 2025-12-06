package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dasciam/bedrockscanner/data"
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
	"time"
)

func main() {
	var (
		dbPath string

		what       string
		pps        int
		numSockets int

		timeout int

		rescan  bool
		rewrite bool
	)

	flag.StringVar(&what, "what", "ALL", "What to scan (subnet, file path or ALL)")
	flag.IntVar(&pps, "packets-per-second", 5_000, "Number of max packers per second")
	flag.StringVar(&dbPath, "db", "result.db", "Path to output/input DB")
	flag.IntVar(&numSockets, "num-sockets", 1, "Number of sockets")
	flag.BoolVar(&rescan, "rescan", false, "Rescan all servers in the database")
	flag.BoolVar(&rewrite, "rewrite", false, "Mark all servers as offline in the database before the scan")
	flag.IntVar(&timeout, "timeout", 60, "Timeout time for sockets (in seconds)")
	flag.Parse()

	if numSockets < 0 {
		panic("num-sockets must be >= 0")
	} else if numSockets == 0 {
		numSockets = 1
	}
	if dbPath == "" {
		panic("db path not specified")
	}
	db := data.New(dbPath)

	log.Printf("Settings:\n- Subnet/file: %s\n- PPS (Packets/s): %d\n- DB: %s", what, pps, dbPath)

	var rng []ranges.Addr

	if rescan {
		rng = append(rng, scanner.RangeFromInput(db))
	} else {
		parsePrefix, err := netip.ParsePrefix(what)
		if err != nil {
			switch {
			case strings.ToLower(what) == "all":
				const partCount = 64
				const part = math.MaxUint32 / partCount

				for i := 0; i < partCount; i++ {
					if i == partCount-1 {
						rng = append(rng, ranges.NewUInt32(uint32(part*i), math.MaxUint32))
						break
					}
					rng = append(rng, ranges.NewUInt32(uint32(part*i), uint32(part*(i+1))))
				}
			default:
				contents, err := os.ReadFile(what)
				if err != nil {
					log.Fatal(err)
				}
				rng = lo.FilterMap(bytes.Split(contents, []byte("\n")), func(v []byte, _ int) (ranges.Addr, bool) {
					prefix, err := netip.ParsePrefix(string(v))
					if err != nil {
						log.Printf("Error parsing prefix from file %s: %v", what, err)
						return nil, false
					}
					return ranges.NewNetIP(prefix), true
				})
			}
		} else {
			rng = append(rng, ranges.NewNetIP(parsePrefix))
		}
	}

	if rewrite {
		db.FlagOffline()
	}

	var sockets []net.PacketConn

	for range numSockets {
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			panic(err)
		}
		sockets = append(sockets, conn)
	}

	outputs := []scanner.Output{
		output.Print{},
		db,
	}

	var (
		wg           sync.WaitGroup
		readWorkerWg sync.WaitGroup
	)

	done := make(chan struct{})

	for _, socket := range sockets {
		readWorkerWg.Add(1)
		go scanner.ReadWorker(socket, output.NewMulti(outputs...), &readWorkerWg, time.Duration(timeout)*time.Second)
	}

	limiter := limit.NewBasicLimiter(pps)

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
