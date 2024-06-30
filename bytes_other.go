//go:build !go1.20
// +build !go1.20

package mem

import (
	"reflect"
	"unsafe"
)

func (r RO) bytes() []byte {
	s := r.str()
	return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}
