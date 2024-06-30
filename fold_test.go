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

package mem // import "go4.org/mem"

import "testing"

type foldTest struct {
	a, b string
	want bool
}

func TestContainsFold(t *testing.T) {
	runFoldTests(t, ContainsFold, "ContainsFold", []foldTest{
		{"foo", "", true},
		{"", "", true},
		{"", "foo", false},
		{"foo", "foo", true},
		{"FOO", "foo", true},
		{"foo", "FOO", true},
		{"foo ", "FOO", true},
		{" foo", "FOO", true},
		{" foo ", "FOO", true},
		{"FOO ", "foo", true},
		{" FOO", "foo", true},
		{" FOO ", "foo", true},
		{" FOO ", "bar", false},
	})
}

func TestHasPrefixFold(t *testing.T) {
	runFoldTests(t, HasPrefixFold, "HasPrefixFold", []foldTest{
		{"foo", "", true},
		{"", "", true},
		{"", "foo", false},
		{"foo", "foo", true},
		{"FOO", "foo", true},
		{"foo", "FOO", true},
		{"foo", "food", false},
	})
}

func TestHasSuffixFold(t *testing.T) {
	runFoldTests(t, HasSuffixFold, "HasSuffixFold", []foldTest{
		{"foo", "", true},
		{"", "", true},
		{"", "foo", false},
		{"foo", "foo", true},
		{"FOO", "foo", true},
		{"foo", "FOO", true},
		{" foo", "FOO", true},
		{" foo", "FoO", true},
	})
}

func runFoldTests(t *testing.T, fn func(RO, RO) bool, funcName string, tests []foldTest) {
	t.Helper()
	for _, tt := range tests {
		got := fn(S(tt.a), S(tt.b))
		if got != tt.want {
			t.Errorf("%s(%q, %q) = %v; want %v", funcName, tt.a, tt.b, got, tt.want)
		}
	}
}
