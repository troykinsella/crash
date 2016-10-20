package util

import "errors"

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type rangeIterator struct {
	i    int
	from int
	to   int
	step int
}

// to: inclusive
func NewRangeIterator(from int, to int, step int) (Iterator, error) {
    if step == 0 {
	    return nil, errors.New("cannot step 0")
    }
	if step > 0 {
		if from > to {
			return nil, errors.New("from > to")
		}
	} else {
		if from < to {
			return nil, errors.New("from < to")
		}
	}

	return &rangeIterator{from, from, to, step}, nil
}

func (ri *rangeIterator) HasNext() bool {
	if ri.step > 0 {
		return ri.i <= ri.to
	}
	return ri.i >= ri.to
}

func (ri *rangeIterator) Next() interface{} {
	if !ri.HasNext() {
		return nil
	}

	i := ri.i
	ri.i += ri.step
	return i
}

type sliceIterator struct {
	i int
	l int
	s []interface{}
}

func NewSliceIterator(s []interface{}) (Iterator, error) {
	if s == nil {
		return nil, errors.New("nil slice")
	}
	return &sliceIterator{0, len(s), s}, nil
}

func (si *sliceIterator) HasNext() bool {
	return si.i < si.l
}

func (si *sliceIterator) Next() interface{} {
	if !si.HasNext() {
		return nil
	}

	i := si.i
	si.i += 1
	return si.s[i]
}
