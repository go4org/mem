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

// Tests copied from Go's strings package. (BSD license)

package mem

import (
	"testing"
	"unicode"
)

var faces = "☺☻☹"

type FieldsTest struct {
	s string
	a []string
}

var fieldstests = []FieldsTest{
	{"", []string{}},
	{" ", []string{}},
	{" \t ", []string{}},
	{"\u2000", []string{}},
	{"  abc  ", []string{"abc"}},
	{"1 2 3 4", []string{"1", "2", "3", "4"}},
	{"1  2  3  4", []string{"1", "2", "3", "4"}},
	{"1\t\t2\t\t3\t4", []string{"1", "2", "3", "4"}},
	{"1\u20002\u20013\u20024", []string{"1", "2", "3", "4"}},
	{"\u2000\u2001\u2002", []string{}},
	{"\n™\t™\n", []string{"™", "™"}},
	{"\n\u20001™2\u2000 \u2001 ™", []string{"1™2", "™"}},
	{"\n1\uFFFD \uFFFD2\u20003\uFFFD4", []string{"1\uFFFD", "\uFFFD2", "3\uFFFD4"}},
	{"1\xFF\u2000\xFF2\xFF \xFF", []string{"1\xFF", "\xFF2\xFF", "\xFF"}},
	{faces, []string{faces}},
}

func eq(a []RO, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].EqualString(b[i]) {
			return false
		}
	}
	return true
}

func TestFields(t *testing.T) {
	for _, tt := range fieldstests {
		a := AppendFields(nil, S(tt.s))
		if !eq(a, tt.a) {
			t.Errorf("Fields(%q) = %v; want %v", tt.s, a, tt.a)
			continue
		}
	}
}

var FieldsFuncTests = []FieldsTest{
	{"", []string{}},
	{"XX", []string{}},
	{"XXhiXXX", []string{"hi"}},
	{"aXXbXXXcX", []string{"a", "b", "c"}},
}

func TestFieldsFunc(t *testing.T) {
	for _, tt := range fieldstests {
		a := AppendFieldsFunc(nil, S(tt.s), unicode.IsSpace)
		if !eq(a, tt.a) {
			t.Errorf("FieldsFunc(%q, unicode.IsSpace) = %v; want %v", tt.s, a, tt.a)
			continue
		}
	}
	pred := func(c rune) bool { return c == 'X' }
	for _, tt := range FieldsFuncTests {
		a := AppendFieldsFunc(nil, S(tt.s), pred)
		if !eq(a, tt.a) {
			t.Errorf("FieldsFunc(%q) = %v, want %v", tt.s, a, tt.a)
		}
	}
}

func TestFieldsAllocs(t *testing.T) {
	var f []RO
	n := int(testing.AllocsPerRun(1000, func() {
		f = AppendFields(f[:0], S(" foo bar baz"))
		if len(f) != 3 {
			panic("wrong result")
		}
	}))
	if n != 0 {
		t.Fatalf("allocs = %d; want 0", n)
	}
}
