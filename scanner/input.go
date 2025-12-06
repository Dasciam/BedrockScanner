package scanner

import (
	"github.com/dasciam/bedrockscanner/ranges"
	"iter"
	"math"
	"net"
	"net/netip"
)

type inputRange struct {
	current  string
	leftover []string
}

func (i inputRange) Addr() netip.Addr {
	host, _, _ := net.SplitHostPort(i.current)
	return netip.MustParseAddr(host)
}

func (i inputRange) Next() (ranges.Addr, bool) {
	if len(i.leftover) <= 0 {
		return nil, false
	}
	return inputRange{i.leftover[0], i.leftover[1:]}, true
}

type Input interface {
	Read() iter.Seq2[string, Pong] // addr, Pong
}

func RangeFromInput(input Input) ranges.Addr {
	var addresses []string
	for addr := range input.Read() {
		addresses = append(addresses, addr)
	}
	if len(addresses) == 0 {
		// No addresses
		return ranges.NewUInt32(0, math.MaxUint32)
	}
	return inputRange{addresses[0], addresses[1:]}
}
