package qb_shamir

import (
	"fmt"
	"io"
)

type reader struct {
	io.Reader
	readers map[byte]io.Reader
	eof     bool
}

func NewReader(readers map[byte]io.Reader) (io.Reader, error) {
	// Verify enough parts provided
	if len(readers) < 2 {
		return nil, fmt.Errorf("at least two parts are required to reconstruct the secret")
	}
	return &reader{readers: readers}, nil
}

func (r *reader) Read(p []byte) (int, error) {
	if r.eof {
		return 0, io.EOF
	}

	points := make([][]pair, len(p))
	for i := range points {
		points[i] = make([]pair, len(r.readers))
	}

	j := 0
	n := 0

	for x, ir := range r.readers {
		buf := make([]byte, len(p))
		m, err := ir.Read(buf)
		if io.EOF == err {
			r.eof = true
		} else if nil != err {
			return 0, err
		} else if 0 != n && 0 != m && m != n {
			return 0, fmt.Errorf("input must be of equal length")
		}
		n = m

		for i := 0; i < m; i++ {
			points[i][j] = pair{x: x, y: buf[i]}
		}
		j++
	}

	for m := 0; m < n; m++ {
		p[m] = interpolate(points[m], 0)
	}

	return n, nil
}
