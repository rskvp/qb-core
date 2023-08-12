package qb_shamir

import (
	"crypto/rand"
	"fmt"
	"io"
)

type writer struct {
	io.Writer
	writers      map[byte]io.Writer
	threshold    int
	bytesWritten int
}

func (w *writer) Write(p []byte) (int, error) {
	n := 0
	// Construct a random polynomial for each byte of the secret.
	// Because we are using a field of size 256, we can only represent
	// a single byte as the intercept of the polynomial, so we must
	// use a new polynomial for each byte.
	for _, val := range p {
		p, err := makePolynomial(val, uint8(w.threshold-1))
		if nil != err {
			return n, fmt.Errorf("failed to generate polynomial: %v", err)
		}

		// Generate a `parts` number of (x,y) pairs
		// We cheat by encoding the x value once as the final index,
		// so that it only needs to be stored once.
		for x, w := range w.writers {
			y := p.evaluate(uint8(x))
			_, err := w.Write([]byte{y})
			if nil != err {
				return n, fmt.Errorf("failed to write part: %v", err)
			}
		}
		n++
		w.bytesWritten += n
	}

	return n, nil
}

func NewWriter(parts, threshold int, factory func(x byte) (io.Writer, error)) (io.Writer, error) {
	// Sanity check the input
	if parts < threshold {
		return nil, fmt.Errorf("parts cannot be less than threshold")
	}
	if parts > 255 {
		return nil, fmt.Errorf("parts cannot exceed 255")
	}
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}

	result := writer{writers: make(map[byte]io.Writer, parts), threshold: threshold}

	buf := make([]byte, 1)
	for len(result.writers) < parts {
		if _, err := rand.Read(buf); err != nil {
			return nil, err
		}
		if _, exists := result.writers[buf[0]]; exists {
			continue
		}
		w, err := factory(buf[0])
		if nil != err {
			return nil, err
		}
		result.writers[buf[0]] = w
	}

	return &result, nil
}
