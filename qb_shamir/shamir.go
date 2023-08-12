package qb_shamir

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

type ShamirHelper struct {
}

var Shamir *ShamirHelper

func init() {
	Shamir = new(ShamirHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShamirHelper) SplitToFiles(secret []byte, parts, threshold int, dir string) ([]string, error) {
	response := make([]string, 0)

	// split parts
	data, err := instance.Split(secret, parts, threshold)
	if nil != err {
		return nil, err
	}

	// save parts to files
	if len(dir) == 0 {
		dir = "./"
	}
	dir = qb_utils.Paths.Absolute(dir)
	_ = qb_utils.Paths.Mkdir(dir + qb_utils.OS_PATH_SEPARATOR)
	name := qb_utils.Coding.MD5(string(secret))
	count := 0
	for k, v := range data {
		filename := qb_utils.Paths.Concat(dir, fmt.Sprintf("%s-%v.part", name, count))
		count++
		_, e := qb_utils.IO.WriteTextToFile(fmt.Sprintf("%v\n%v", k, qb_utils.Coding.EncodeBase64(v)), filename)
		if nil != e {
			return nil, e
		}
		response = append(response, filename)
	}

	return response, nil
}

// Split takes an arbitrarily long secret and generates a `parts`
// number of shares, `threshold` of which are required to reconstruct
// the secret. The parts and threshold must be at least 2, and less
// than 256. The returned shares are each one byte longer than the secret
// as they attach a tag used to reconstruct the secret.
func (instance *ShamirHelper) Split(secret []byte, parts, threshold int) (map[byte][]byte, error) {
	buffers := make(map[byte]*bytes.Buffer, parts)
	factory := func(x byte) (io.Writer, error) {
		buffers[x] = &bytes.Buffer{}
		return buffers[x], nil
	}
	s, err := NewWriter(parts, threshold, factory)
	if nil != err {
		return nil, fmt.Errorf("failed to initilize writer: %v", err)
	}

	if _, err := s.Write(secret); nil != err {
		return nil, fmt.Errorf("failed to split secret: %v", err)
	}

	out := make(map[byte][]byte, parts)
	for x, buf := range buffers {
		out[x] = buf.Bytes()
	}

	// Return the encoded secrets
	return out, nil
}

func (instance *ShamirHelper) CombineFromDir(dir string) ([]byte, error) {
	files, err := qb_utils.Paths.ListFiles(dir, "*.part")
	if nil != err {
		return nil, err
	}
	return instance.CombineFromFiles(files)
}

func (instance *ShamirHelper) CombineFromFiles(files []string) ([]byte, error) {
	parts := make(map[byte][]byte)
	for _, filename := range files {
		text, err := qb_utils.IO.ReadTextFromFile(filename)
		if nil != err {
			return nil, err
		}
		tokens := strings.Split(text, "\n")
		if len(tokens) > 1 {
			k := qb_utils.Convert.ToInt(tokens[0])
			v, e := qb_utils.Coding.DecodeBase64(strings.Join(tokens[1:], ""))
			if nil != e {
				return nil, e
			}

			parts[byte(k)] = v
		} else {
			return nil, errors.New(fmt.Sprintf("invalid file size or format: %s", filename))
		}
	}

	return instance.Combine(parts)
}

// Combine is used to reverse a Split and reconstruct a secret
// once a `threshold` number of parts are available.
func (instance *ShamirHelper) Combine(parts map[byte][]byte) ([]byte, error) {
	// Verify enough parts provided
	if len(parts) < 2 {
		return nil, fmt.Errorf("less than two parts cannot be used to reconstruct the secret")
	}

	// Verify the parts are all the same length
	var firstPartLen int
	for x := range parts {
		firstPartLen = len(parts[x])
		break
	}
	if firstPartLen < 1 {
		return nil, fmt.Errorf("parts must be at least one byte long")
	}
	for _, part := range parts {
		if len(part) != firstPartLen {
			return nil, fmt.Errorf("all parts must be the same length")
		}
	}

	// Create a buffer to store the reconstructed secret
	secret := make([]byte, firstPartLen)
	points := make([]pair, len(parts))

	for i := range secret {
		p := 0
		for k, v := range parts {
			points[p] = pair{x: k, y: v[i]}
			p++
		}
		secret[i] = interpolate(points, 0)
	}

	return secret, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// an x/y pair
type pair struct {
	x, y byte
}

// polynomial represents a polynomial of arbitrary degree
type polynomial struct {
	coefficients []uint8
}

// makePolynomial constructs a random polynomial of the given
// degree but with the provided intercept value.
func makePolynomial(intercept, degree uint8) (polynomial, error) {
	// Create a wrapper
	p := polynomial{
		coefficients: make([]byte, degree+1),
	}

	// Ensure the intercept is set
	p.coefficients[0] = intercept

	// Assign random co-efficients to the polynomial
	if _, err := rand.Read(p.coefficients[1:]); err != nil {
		return p, err
	}

	return p, nil
}

// evaluate returns the value of the polynomial for the given x
func (p *polynomial) evaluate(x byte) byte {
	// Special case the origin
	if x == 0 {
		return p.coefficients[0]
	}

	// Compute the polynomial value using Horner's method.
	degree := len(p.coefficients) - 1
	out := p.coefficients[degree]
	for i := degree - 1; i >= 0; i-- {
		coeff := p.coefficients[i]
		out = add(mult(out, x), coeff)
	}
	return out
}

// Lagrange interpolation
//
// Takes N sample points and returns the value at a given x using a lagrange interpolation.
func interpolate(points []pair, x byte) (value byte) {
	for i, a := range points {
		weight := byte(1)
		for j, b := range points {
			if i != j {
				top := x ^ b.x
				bottom := a.x ^ b.x
				factor := div(top, bottom)
				weight = mult(weight, factor)
			}
		}
		value = value ^ mult(weight, a.y)
	}
	return
}

// div divides two numbers in GF(2^8)
func div(a, b uint8) uint8 {
	if b == 0 {
		// leaks some timing information but we don't care anyways as this
		// should never happen, hence the panic
		panic("divide by zero")
	}

	var goodVal, zero uint8
	log_a := logTable[a]
	log_b := logTable[b]
	diff := (int(log_a) - int(log_b)) % 255
	if diff < 0 {
		diff += 255
	}

	ret := expTable[diff]

	// Ensure we return zero if a is zero but aren't subject to timing attacks
	goodVal = ret

	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}

	return ret
}

// mult multiplies two numbers in GF(2^8)
func mult(a, b uint8) (out uint8) {
	var goodVal, zero uint8
	log_a := logTable[a]
	log_b := logTable[b]
	sum := (int(log_a) + int(log_b)) % 255

	ret := expTable[sum]

	// Ensure we return zero if either a or be are zero but aren't subject to
	// timing attacks
	goodVal = ret

	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}

	if subtle.ConstantTimeByteEq(b, 0) == 1 {
		ret = zero
	} else {
		// This operation does not do anything logically useful. It
		// only ensures a constant number of assignments to thwart
		// timing attacks.
		goodVal = zero
	}

	return ret
}

// add combines two numbers in GF(2^8)
// This can also be used for subtraction since it is symmetric.
func add(a, b uint8) uint8 {
	return a ^ b
}
