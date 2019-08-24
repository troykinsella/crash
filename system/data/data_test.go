package data

import (
	"testing"
)

func TestLooseEqual(t *testing.T) {
	var tests = []struct {
		l interface{}
		r interface{}
		eq bool
		//err string
	}{
		// nils
		{nil,   nil,   true},
		{nil,   "",    true},
		{"",    nil,   true},
		{"",    "",    true},
		{nil,   0,     false},
		{0,     nil,   false},
		{"",    0,     false},
		{0,     "",    false},
		{true,  nil,   false},
		{nil,   true,  false},
		{false, nil,   false},
		{nil,   false, false},
		{"foo", nil,   false},
		{nil,   "foo", false},
		{"foo", "",    false},
		{"",    "foo", false},

		// bools
		{true,    true,    true},
		{true,    false,   false},
		{false,   true,    false},
		{false,   false,   true},
		{"foo",   true,    false},
		{true,    "foo",   false},
		{"foo",   false,   false},
		{false,   "foo",   false},
		{"true",  true,    true},
		{true,    "true",  true},
		{"true",  "true",  true},
		{"true",  false,   false},
		{false,   "true",  false},
		{"false", true,    false},
		{true,    "false", false},
		{"true",  false,   false},
		{"false", false,   true},
		{false,   "false", true},
		{"false", "false", true},

		// strings
		{"foo", "bar",       false},
		{"bar", "foo",       false},
		{"foo", "foo",       true},
		{"foo\n", "bar\n",   false},
		{"foo\n", "foo\n",   true},
		{"\r\n\t", "\n\r\t", false},
		{"\r\n\t", "\r\n\t", true},

		// integers
		{0,  0,  true},
		{1,  0,  false},
		{0,  1,  false},
		{1,  1,  true},
		{-1, 0,  false},
		{0,  -1, false},
		{-1, -1, true},

		// mixed strings / integers
		{"foo", 0,     false},
		{0,     "foo", false},
		{"0",   0,     true},
		{0,     "0",   true},
		{"1",   1,     true},
		{1,     "1",   true},
		{"-1",  -1,    true},
		{-1,    "-1",  true},
	}

	for _, test := range tests {
		eq, err := LooseEqual(test.l, test.r)
		if err != nil {
			t.Errorf("error: %s", err.Error())
		} else if eq != test.eq {
			t.Errorf("%v. unexpected equality: expected=%t actual=%t", test, test.eq, eq)
		}
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		l interface{}
		r interface{}
		cmp int
		//err string
	}{
		// nils
		{nil,   nil,   0},
		{"",    "",    0},
		{nil,   "",    0},
		{"",    nil,   0},

		// strings
		{"a",  "a",  0},
		{"a",  "b",  -1},
		{"b",  "a",  1},

		// integers
		{0,   0,   0},
		{1,   0,   1},
		{0,   1,   -1},

		// bools
		{true,  true,  0},
		{false, false, 0},
		{true,  false, 1},
		{false, true,  -1},

		// mixed strings / integers
		{"0",  0,    0},
		{0,    "0",  0},
		{"1",  0,    1},
		{1,    "0",  1},
		{"0",  1,    -1},
		{0,    "1",  -1},
		{"-1", 0,    -1},
		{-1,   "0",  -1},
	}

	for _, test := range tests {
		cmp, err := Compare(test.l, test.r)
		if err != nil {
			t.Errorf("error: %s", err.Error())
		} else if cmp != test.cmp {
			t.Errorf("%v. unexpected compare result: expected=%d actual=%d", test, test.cmp, cmp)
		}
	}
}
