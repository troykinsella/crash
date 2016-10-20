package scanner

import (
	"testing"
	"bufio"
	"strings"
)

func assertPos(p int, r *RuneReader, t *testing.T) {
	if r.Pos() != p {
		t.Errorf("Unexpected pos: expected=%d actual=%d", p, r.Pos())
	}
}

func runeEq(actual rune, expected rune, t *testing.T) {
	if actual != expected {
		t.Errorf("Unexpected rune: expected=%c actual=%c", expected, actual)
	}
}

func TestRuneReader_all(t *testing.T) {
	str := "abcdefg"
	in := bufio.NewReader(strings.NewReader(str))
	r := NewRuneReader(in)
	assertPos(0, r, t)

	ch, _ := r.Read()
	runeEq(ch, 'a', t)
	assertPos(1, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'b', t)
	assertPos(2, r, t)

	r.Unread()

	ch, _ = r.Read()
	runeEq(ch, 'b', t)
	assertPos(2, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'c', t)
	assertPos(3, r, t)

	r.Rewind()

	ch, _ = r.Read()
	runeEq(ch, 'a', t)
	assertPos(1, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'b', t)
	assertPos(2, r, t)

	r.Unread()

	ch, _ = r.Read()
	runeEq(ch, 'b', t)
	assertPos(2, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'c', t)
	assertPos(3, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'd', t)
	assertPos(4, r, t)

	r.Reset()

	ch, _ = r.Read()
	runeEq(ch, 'e', t)
	assertPos(5, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'f', t)
	assertPos(6, r, t)

	r.Rewind()

	ch, _ = r.Read()
	runeEq(ch, 'e', t)
	assertPos(5, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'f', t)
	assertPos(6, r, t)

	ch, _ = r.Read()
	runeEq(ch, 'g', t)
	assertPos(7, r, t)
}
