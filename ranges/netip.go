package ranges

import "net/netip"

type NetIPAddr struct {
	prefix netip.Prefix
	addr   netip.Addr
}

func NewNetIP(prefix netip.Prefix) NetIPAddr {
	return NetIPAddr{
		prefix: prefix,
		addr:   prefix.Addr(),
	}
}

func (n NetIPAddr) Addr() netip.Addr {
	return n.addr
}

func (n NetIPAddr) Next() (Addr, bool) {
	next := n.addr.Next()
	if !n.prefix.Contains(next) {
		return nil, false
	}
	return NetIPAddr{
		prefix: n.prefix,
		addr:   next,
	}, true
}
