package system

import (
	"os"
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/system/parser"
)

func Interpolate(str string, values util.Values) string {
	str = os.Expand(str, func(match string) string {
		val := values.GetString(match)
		return val
	})
	return str
}

type Interpolator struct {
	i *Interpreter
	p *parser.Parser
}

func NewInterpolator(i *Interpreter, p *parser.Parser) *Interpolator {
	return &Interpolator{
		i: i,
		p: p,
	}
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
	return Interpolate(val, i.vars)
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
