package system

import (
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/system/parser"
	"bufio"
	"strings"
	"github.com/troykinsella/crash/system/exec"
)

func Interpolate(str string, vars util.Values) (string, error) {
	r := bufio.NewReader(strings.NewReader(str))
	p := parser.New(r)

	istr, err := p.IString()
	if err != nil {
		return "", err
	}

	ctx := exec.NewContext(vars, nil)
	_, result, err := istr.Exec(ctx)
	if err != nil {
		return "", err
	}

	return result, nil
}

type interpVals struct {
	wrapped util.Values
	vars util.Values
}

func (i *interpVals) Get(name string) interface{} {
	return i.GetString(name)
}

func (i *interpVals) GetString(name string) string {
	val := i.wrapped.GetString(name)
	str, err := Interpolate(val, i.vars)
	if err != nil {
		panic(err)
	}
	return str
}

func (i *interpVals) AsMap() map[string]interface{} {
	// TODO: interp each value
	return i.wrapped.AsMap()
}

func NewInterpolatedValues(wrapped util.Values, vars util.Values) util.Values {
	if vars == nil {
		vars = wrapped
	}

	return &interpVals{
		wrapped: wrapped,
		vars: vars,
	}
}
