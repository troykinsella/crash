package util

type Bool struct {
	b bool
}

func (b *Bool) Value() bool {
	return b.b
}

var True = &Bool{true}
var False = &Bool{false}

func BoolFor(val bool) *Bool {
	if val {
		return True
	}
	return False
}
