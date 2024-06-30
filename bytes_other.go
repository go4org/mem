//go:build !go1.20
// +build !go1.20

package mem

import (
	"reflect"
	"unsafe"
)

func (r RO) bytes() []byte {
	s := r.str()
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  len(s),
		Cap:  len(s),
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
	}))
}
