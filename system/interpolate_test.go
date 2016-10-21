package system

import (
	"testing"
	"github.com/troykinsella/crash/util"
)

func TestInterpolate(t *testing.T) {
	vars := util.AsValues(map[string]interface{}{
		"foo": "few",
		"bar": "bahr",
	})

	var tests = []struct {
		fmt string
		out string
		err string
	}{
		{ "foo", "foo", "" },
		{ "$foo", "few", "" },
		{ "${foo}", "few", "" },
		{ "bar$foo", "barfew", "" },
		{ "bar${foo}", "barfew", "" },
		{ "bar${foo}baz", "barfewbaz", "" },
		{ "$foo$bar", "fewbahr", "" },
		{ "${foo}$bar", "fewbahr", "" },
		{ "$foo${bar}", "fewbahr", "" },
		{ "${foo}${bar}", "fewbahr", "" },
		{ "a${foo}${bar}b", "afewbahrb", "" },
		{ "a${foo}b${bar}c", "afewbbahrc", "" },
		{ "a$foo.stuff", "afew.stuff", "" },

		{ "$foobar", "", "Not found: foobar" },
		{ "${foobar}", "", "Not found: foobar" },
		{ "bar$foobar", "", "Not found: foobar" },
		{ "bar${foobar}", "", "Not found: foobar" },
		{ "bar${foobar}baz", "", "Not found: foobar" },
	}

	for i, test := range tests {
		r, err := Interpolate(test.fmt, vars)
		if test.err == "" {
			if err != nil {
				t.Errorf("%d. \"%s\" unexpected error: %s\n", i, test.fmt, err.Error())
			} else if r != test.out {
				t.Errorf("%d. \"%s\" unexpected result:\nexpected=%s,\nactual=%s\n", i, test.fmt, test.out, r)
			}
		} else {
			if err == nil {
				t.Errorf("%d. \"%s\" expected error:\nexpected=%s,\nactual=nil\n", i, test.fmt, test.err)
			} else if test.err != err.Error() {
				t.Errorf("%d. \"%s\" unexpected error:\nexpected=%s,\nactual=%s\n", i, test.fmt, test.err, err.Error())
			}
		}
	}
}
