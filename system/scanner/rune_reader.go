package scanner

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

type RuneReader struct {
	r *bufio.Reader
	pos int

	rewound bool
	buf *bytes.Buffer
}

func NewRuneReader(r *bufio.Reader) *RuneReader {
	return &RuneReader{
		r: r,
		buf: &bytes.Buffer{},
	}
}

func (rr *RuneReader) Pos() int {
	return rr.pos
}

func (rr *RuneReader) Read() (rune, error) {
	if rr.rewound {
		ch, _, err := rr.buf.ReadRune()
		if err != io.EOF {
			rr.pos += 1
			return ch, nil
		}
		rr.rewound = false
	}

	ch, _, err := rr.r.ReadRune()
	if err != nil {
		return -1, err
	}

	rr.pos += 1
	rr.buf.WriteRune(ch)

	return ch, nil
}

func (rr *RuneReader) Unread() error {
	if rr.rewound {
		err := rr.buf.UnreadRune()
		if err != nil {
			return err
		}
	} else {
		l := rr.buf.Len()
		if l == 0 {
			return errors.New("cannot unread at position 0")
		}

		err := rr.r.UnreadRune()
		if err != nil {
			return err
		}
		rr.buf.Truncate(l - 1)
	}

	rr.pos -= 1

	return nil
}

func (rr *RuneReader) Rewind() error {
	if rr.rewound {
		return errors.New("Cannot rewind twice")
	}
	rr.rewound = true
	rr.pos -= rr.buf.Len()

	return nil
}

func (rr *RuneReader) Reset() {
	rr.rewound = false
	rr.buf.Reset()
}
