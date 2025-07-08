package ranges

import (
	"net/netip"
)

type Addr interface {
	Addr() netip.Addr
	Next() (Addr, bool)
}
