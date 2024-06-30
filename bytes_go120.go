//go:build go1.20
// +build go1.20

package mem

import "unsafe"

// get a unsafe bytes view of our unsafe string
func (r RO) bytes() []byte {
	s := r.str()
	d := unsafe.StringData(s)
	return unsafe.Slice(d, len(s))
}
