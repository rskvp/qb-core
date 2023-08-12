package qb_rnd

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

type RndHelper struct {
}

var Rnd *RndHelper

func init() {
	Rnd = new(RndHelper)
}

var (
	NUMBERS     = "1234567890"
	CHARSET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	CHARSET_LOW = "0123456789abcdefghijklmnopqrstuvwxyz"
	CHARSET_UP  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var uuidRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var uuidMutex = &sync.Mutex{}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *RndHelper) NewValuesRandomizer(args ...interface{}) (*ValuesRandomizer, error) {
	return NewValuesRandomizer(args...)
}

// Uuid creates a new random UUID or panics.
func (instance *RndHelper) Uuid() string {
	return uuid.New().String()
}

func (instance *RndHelper) RndId() string {
	return randomString(32, CHARSET)
}

func (instance *RndHelper) UuidTimestamp() string {
	return time.Now().Format("20060102T150405") + "-" + instance.Uuid()
}

func (instance *RndHelper) Between(min, max int64) int64 {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func (instance *RndHelper) BetweenDuration(max, min int64) time.Duration {
	if min == max {
		return time.Duration(min)
	}
	return time.Duration(instance.Between(max, min))
}

func (instance *RndHelper) RndDigits(n int) string {
	return randomString(n, NUMBERS)
}

func (instance *RndHelper) RndChars(n int) string {
	return randomString(n, CHARSET)
}

func (instance *RndHelper) RndCharsLower(n int) string {
	return randomString(n, CHARSET_LOW)
}

func (instance *RndHelper) RndCharsUpper(n int) string {
	return randomString(n, CHARSET_UP)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func randomString(l int, pool string) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// uuidRFC4122 generates a random UUID according to RFC 4122.
func uuidRFC4122() string {
	uuidArray := make([]byte, 16)
	uuidMutex.Lock()
	_, _ = uuidRand.Read(uuidArray)
	uuidMutex.Unlock()
	// variant bits; see section 4.1.1
	uuidArray[8] = uuidArray[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuidArray[6] = uuidArray[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidArray[0:4], uuidArray[4:6], uuidArray[6:8], uuidArray[8:10], uuidArray[10:])
}
