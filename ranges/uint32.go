package ranges

import (
	"math"
	"net/netip"
	"slices"
	"unsafe"
)

type UInt32 struct {
	v uint32

	constraint uint32
}

func NewUInt32(v uint32, constraint uint32) *UInt32 {
	return &UInt32{v: v, constraint: constraint}
}

func (u UInt32) Addr() netip.Addr {
	x := *(*[4]byte)(unsafe.Pointer(&u.v))
	slices.Reverse(x[:])
	return netip.AddrFrom4(x)
}

func (u UInt32) Next() (Addr, bool) {
	next := u.v + 1
	if next >= math.MaxUint32 || next >= u.constraint {
		return nil, false
	}
	return UInt32{v: next, constraint: u.constraint}, true
}
