package qb_vcal

import (
	"bufio"
	"io"
)

type CalendarStream struct {
	r io.Reader
	b *bufio.Reader
}

func NewCalendarStream(r io.Reader) *CalendarStream {
	return &CalendarStream{
		r: r,
		b: bufio.NewReader(r),
	}
}

func (instance *CalendarStream) ReadLine() (*ContentLine, error) {
	response := make([]byte, 0)
	loop := true
	var err error
	for loop {
		var b []byte
		b, err = instance.b.ReadBytes('\n')
		if len(b) == 0 {
			if err == nil {
				continue
			} else {
				loop = false
			}
		} else if b[len(b)-1] == '\n' {
			o := 1
			if len(b) > 1 && b[len(b)-2] == '\r' {
				o = 2
			}
			p, err := instance.b.Peek(1)
			response = append(response, b[:len(b)-o]...)
			if err == io.EOF {
				loop = false
			}
			if len(p) == 0 {
				loop = false
			} else if p[0] == ' ' {
				_, _ = instance.b.Discard(1)
			} else {
				loop = false
			}
		} else {
			response = append(response, b...)
		}
		switch err {
		case nil:
			if len(response) == 0 {
				loop = true
			}
		case io.EOF:
			loop = false
		default:
			return nil, err
		}
	}
	if len(response) == 0 && err != nil {
		return nil, err
	}
	cl := ContentLine(response)
	return &cl, err
}
