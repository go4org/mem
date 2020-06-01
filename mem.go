/*
Copyright 2020 The Go4 AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mem provides the mem.RO type that allows you to cheaply pass &
// access either a read-only []byte or a string.
package mem // import "go4.org/mem"

import (
	"strconv"
	"strings"
	"unsafe"
)

// RO is a read-only view of some bytes of memory. It may be be backed
// by a string or []byte. Notably, unlike a string, the memory is not
// guaranteed to be immutable. While the length is fixed, the
// underlying bytes might change if interleaved with code that's
// modifying the underlying memory.
//
// RO is a value type that's the same size of a Go string. Its various
// methods should inline & compile to the equivalent operations
// working on a string or []byte directly.
//
// Unlike a Go string, RO is not 'comparable' (it can't be a map key
// or support ==). Use its Equal method to compare. This is done so an
// RO backed by a later-mutating []byte doesn't break invariants in
// Go's map implementation.
type RO struct {
	_ [0]func() // not comparable; don't want to be a map key or support ==
	m unsafeString
}

// str returns the unsafeString as a string. Only for us with standard
// library funcs known to not let the string escape, as it doesn't
// obey the language/runtime's expectations of a real string (it can
// change underfoot).
func (r RO) str() string { return string(r.m) }

// Len returns len(r).
func (r RO) Len() int { return len(r.m) }

// At returns r[i].
func (r RO) At(i int) byte { return r.m[i] }

// Slice returns r[from:to].
func (r RO) Slice(from, to int) RO { return RO{m: r.m[from:to]} }

// SliceFrom returns r[from:].
func (r RO) SliceFrom(from int) RO { return RO{m: r.m[from:]} }

// SliceTo returns r[:to].
func (r RO) SliceTo(to int) RO { return RO{m: r.m[:to]} }

// Copy copies up to len(dest) bytes into dest from r and returns the
// number of bytes copied, the min(r.Len(), len(dest)).
func (r RO) Copy(dest []byte) int { return copy(dest, r.m) }

// Equal reports whether r and r2 are the same length and contain the
// same bytes.
func (r RO) Equal(r2 RO) bool { return r.m == r2.m }

// EqualString reports whether r and s are the same length and contain
// the same bytes.
func (r RO) EqualString(s string) bool { return r.str() == s }

// EqualBytes reports whether r and b are the same length and contain
// the same bytes.
func (r RO) EqualBytes(b []byte) bool { return r.str() == string(b) }

// ParseInt returns a signed integer from m, using strconv.ParseInt.
func ParseInt(m RO, base, bitSize int) (int64, error) {
	return strconv.ParseInt(m.str(), base, bitSize)
}

// ParseUint returns a unsigned integer from m, using strconv.ParseUint.
func ParseUint(m RO, base, bitSize int) (uint64, error) {
	return strconv.ParseUint(m.str(), base, bitSize)
}

// Append appends m to dest, and returns the possibly-reallocated
// dest.
func Append(dest []byte, m RO) []byte { return append(dest, m.m...) }

// Contains reports whether substr is within m.
func Contains(m, substr RO) bool { return strings.Contains(m.str(), substr.str()) }

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, which is a more general form
// of case-insensitivity.
func EqualFold(m, m2 RO) bool { return strings.EqualFold(m.str(), m2.str()) }

// HasPrefix reports whether m starts with prefix.
func HasPrefix(m, prefix RO) bool { return strings.HasPrefix(m.str(), prefix.str()) }

// HasSuffix reports whether m ends with suffix.
func HasSuffix(m, suffix RO) bool { return strings.HasSuffix(m.str(), suffix.str()) }

// Index returns the index of the first instance of substr in m, or -1
// if substr is not present in m.
func Index(m, substr RO) int { return strings.Index(m.str(), substr.str()) }

// IndexByte returns the index of the first instance of c in m, or -1
// if c is not present in m.
func IndexByte(m RO, c byte) int { return strings.IndexByte(m.str(), c) }

// LastIndexByte returns the index into m of the last Unicode code
// point satisfying f(c), or -1 if none do.
func LastIndexByte(m RO, c byte) int { return strings.LastIndexByte(m.str(), c) }

// LastIndex returns the index of the last instance of substr in m, or
// -1 if substr is not present in m.
func LastIndex(m, substr RO) int { return strings.LastIndex(m.str(), substr.str()) }

// TrimSpace returns a slice of the string s, with all leading and
// trailing white space removed, as defined by Unicode.
func TrimSpace(m RO) RO { return S(strings.TrimSpace(m.str())) }

// TrimSuffix returns m without the provided trailing suffix.
// If m doesn't end with suffix, m is returned unchanged.
func TrimSuffix(m, suffix RO) RO {
	return S(strings.TrimSuffix(m.str(), suffix.str()))
}

// TrimPrefix returns m without the provided leading prefix.
// If m doesn't start with prefix, m is returned unchanged.
func TrimPrefix(m, prefix RO) RO {
	return S(strings.TrimPrefix(m.str(), prefix.str()))
}

// TrimRightCutset returns a slice of m with all trailing Unicode code
// points contained in cutset removed.
//
// To remove a suffix, use TrimSuffix instead.
func TrimRightCutset(m, cutset RO) RO {
	return S(strings.TrimRight(m.str(), cutset.str()))
}

// TrimLeftCutset returns a slice of m with all leading Unicode code
// points contained in cutset removed.
//
// To remove a prefix, use TrimPrefix instead.
func TrimLeftCutset(m, cutset RO) RO {
	return S(strings.TrimLeft(m.str(), cutset.str()))
}

// TrimCutset returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutset removed.
func TrimCutset(m, cutset RO) RO {
	return S(strings.Trim(m.str(), cutset.str()))
}

// TrimFunc returns a slice of m with all leading and trailing Unicode
// code points c satisfying f(c) removed.
func TrimFunc(m RO, f func(rune) bool) RO {
	return S(strings.TrimFunc(m.str(), f))
}

// TrimRightFunc returns a slice of m with all trailing Unicode
// code points c satisfying f(c) removed.
func TrimRightFunc(m RO, f func(rune) bool) RO {
	return S(strings.TrimRightFunc(m.str(), f))
}

// TrimLeftFunc returns a slice of m with all leading Unicode
// code points c satisfying f(c) removed.
func TrimLeftFunc(m RO, f func(rune) bool) RO {
	return S(strings.TrimLeftFunc(m.str(), f))
}

// NewReader returns a new Reader that reads from m.
func NewReader(m RO) *Reader {
	return &Reader{sr: strings.NewReader(m.str())}
}

// Reader is like a bytes.Reader or strings.Reader.
type Reader struct {
	sr *strings.Reader
}

func (r *Reader) Len() int                                     { return r.sr.Len() }
func (r *Reader) Size() int64                                  { return r.sr.Size() }
func (r *Reader) Read(b []byte) (int, error)                   { return r.sr.Read(b) }
func (r *Reader) ReadAt(b []byte, off int64) (int, error)      { return r.sr.ReadAt(b, off) }
func (r *Reader) ReadByte() (byte, error)                      { return r.sr.ReadByte() }
func (r *Reader) ReadRune() (ch rune, size int, err error)     { return r.sr.ReadRune() }
func (r *Reader) Seek(offset int64, whence int) (int64, error) { return r.sr.Seek(offset, whence) }

// TODO: add Reader.WriteTo, but don't use strings.Reader.WriteTo because it uses io.WriteString, leaking our unsafe string

// unsafeString is a string that's not really a Go string.
// It might be pointing into a []byte. Don't let it escape to callers.
// We contain the unsafety to this package.
type unsafeString string

// S returns a read-only view of the string s.
//
// The compiler should compile this call to nothing. Think of it as a
// free type conversion. The returned RO view is the same size as a
// string.
func S(s string) RO { return RO{m: unsafeString(s)} }

// B returns a read-only view of the byte slice b.
//
// The compiler should compile this call to nothing. Think of it as a
// free type conversion. The returned value is actually smaller than a
// []byte though (16 bytes instead of 24 bytes on 64-bit
// architectures).
func B(b []byte) RO {
	if len(b) == 0 {
		return RO{m: ""}
	}
	type stringHeader struct {
		P   *byte
		Len int
	}
	return RO{m: *(*unsafeString)(unsafe.Pointer(&stringHeader{&b[0], len(b)}))}
}
