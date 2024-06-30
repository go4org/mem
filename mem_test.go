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

package mem

import (
	"io"
	"testing"
)

func TestRO(t *testing.T) {
	b := []byte("some memory.")
	s := "some memory."
	rb := B(b)
	rs := S(s)
	if !rb.Equal(rs) {
		t.Fatal("rb != rs")
	}
	if !rb.EqualString(s) {
		t.Errorf("not equal string")
	}
	if !rs.EqualBytes(b) {
		t.Errorf("not equal byte")
	}
	if !rb.EqualBytes(b) {
		t.Errorf("not equal bytes")
	}
	if !rs.EqualString(s) {
		t.Errorf("not equal string")
	}
	if rb.Less(rs) {
		t.Errorf("bad less")
	}
	if rs.Less(rb) {
		t.Errorf("bad less")
	}
	if !rs.Less(S("~")) {
		t.Errorf("bad less")
	}
	if !rb.Less(S("~")) {
		t.Errorf("bad less")
	}

	if rb.At(0) != 's' {
		t.Fatalf("[0] = %q; want 's'", rb.At(0))
	}
	b[0] = 'z'
	if rb.At(0) != 'z' {
		t.Fatalf("[0] = %q; want 'z'", rb.At(0))
	}

	var b2 = []byte("b2")
	bb := B(b2)
	s = bb.StringCopy()
	b2[0] = '0'
	if s != "b2" {
		t.Fatalf(".StringCopy() is not actually copy, got %q; want %q", s, "b2")
	}

	var got []byte
	got = Append(got, rb)
	got = Append(got, rs)
	want := "zome memory.some memory."
	if string(got) != want {
		t.Errorf("got %q; want %q", got, want)
	}

}

func TestAllocs(t *testing.T) {
	b := []byte("some memory.")
	n := uint(testing.AllocsPerRun(5000, func() {
		ro := B(b)
		if ro.Len() != len(b) {
			t.Fatal("wrong length")
		}
	}))
	if n != 0 {
		t.Errorf("B: unexpected allocs (%d)", n)
	}

	ro := B(b)
	s := string(b)
	n = uint(testing.AllocsPerRun(5000, func() {
		globalString = ro.StringCopy()
		if globalString != s {
			t.Fatal("wrong string")
		}
	}))
	if n != 1 {
		t.Errorf("StringCopy: unexpected allocs (%d)", n)
	}
}

var globalString string

func TestStrconv(t *testing.T) {
	b := []byte("1234")
	i, err := ParseInt(B(b), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1234 {
		t.Errorf("got %d; want 1234", i)
	}
}

var cutTests = []struct {
	s, sep        string
	before, after string
	found         bool
}{
	{"abc", "b", "a", "c", true},
	{"abc", "a", "", "bc", true},
	{"abc", "c", "ab", "", true},
	{"abc", "abc", "", "", true},
	{"abc", "", "", "abc", true},
	{"abc", "d", "abc", "", false},
	{"", "d", "", "", false},
	{"", "", "", "", true},
}

func TestCut(t *testing.T) {
	for _, tt := range cutTests {
		if before, after, found := Cut(S(tt.s), S(tt.sep)); !before.Equal(S(tt.before)) || !after.Equal(S(tt.after)) || found != tt.found {
			t.Errorf("Cut(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before.StringCopy(), after.StringCopy(), found, tt.before, tt.after, tt.found)
		}
	}
}

var cutPrefixTests = []struct {
	s, sep string
	after  string
	found  bool
}{
	{"abc", "a", "bc", true},
	{"abc", "abc", "", true},
	{"abc", "", "abc", true},
	{"abc", "d", "abc", false},
	{"", "d", "", false},
	{"", "", "", true},
}

func TestCutPrefix(t *testing.T) {
	for _, tt := range cutPrefixTests {
		s, sep := S(tt.s), S(tt.sep)
		if after, found := CutPrefix(s, sep); !after.Equal(S(tt.after)) || found != tt.found {
			t.Errorf("CutPrefix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, after.StringCopy(), found, tt.after, tt.found)
		}
	}
}

var cutSuffixTests = []struct {
	s, sep string
	after  string
	found  bool
}{
	{"abc", "bc", "a", true},
	{"abc", "abc", "", true},
	{"abc", "", "abc", true},
	{"abc", "d", "abc", false},
	{"", "d", "", false},
	{"", "", "", true},
}

func TestCutSuffix(t *testing.T) {
	for _, tt := range cutSuffixTests {
		if after, found := CutSuffix(S(tt.s), S(tt.sep)); !after.Equal(S(tt.after)) || found != tt.found {
			t.Errorf("CutSuffix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, after.StringCopy(), found, tt.after, tt.found)
		}
	}
}

func BenchmarkStringCopy(b *testing.B) {
	b.ReportAllocs()
	ro := S("only a fool starts a large fire.")
	for i := 0; i < b.N; i++ {
		globalString = ro.StringCopy()
	}
}

func BenchmarkHash(b *testing.B) {
	b.ReportAllocs()
	ro := S("A man with a beard was always a little suspect anyway.")
	x := ro.MapHash()
	for i := 0; i < b.N; i++ {
		if x != ro.MapHash() {
			b.Fatal("hash changed")
		}
	}
}

// very old go like 1.14 doesn't have io.Discard
type discordWriter struct{}

func (discordWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var discord io.Writer = discordWriter{}

func BenchmarkWriteTo(b *testing.B) {
	b.ReportAllocs()
	ro := S("A man with a beard was always a little suspect anyway.")
	for i := 0; i < b.N; i++ {
		ro.WriteTo(discord)
	}
}

func BenchmarkReader(b *testing.B) {
	b.ReportAllocs()
	ro := S("A man with a beard was always a little suspect anyway.")
	for i := 0; i < b.N; i++ {
		io.Copy(discord, NewReader(ro))
	}
}
