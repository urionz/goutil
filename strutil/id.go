package strutil

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/urionz/goutil/mathutil"
)

var (
	DefMinInt = 1000
	DefMaxInt = 9999
)

// MicroTimeID generate.
// return like: 16074145697981929446(len: 20)
func MicroTimeID() string {
	ms := time.Now().UnixNano() / 1000
	ri := mathutil.RandomInt(DefMinInt, DefMaxInt)

	return strconv.FormatInt(ms, 10) + strconv.FormatInt(int64(ri), 10)
}

// MicroTimeHexID generate.
// return like: 5b5f0588af1761ad3(len: 16-17)
func MicroTimeHexID() string {
	ms := time.Now().UnixNano() / 1000
	ri := mathutil.RandomInt(DefMinInt, DefMaxInt)

	return strconv.FormatInt(ms, 16) + strconv.FormatInt(int64(ri), 16)
}

type UniqIdParams struct {
	Prefix      string
	MoreEntropy bool
}

var entropy = int64(math.Floor(rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * 0x75bcd15))

func NewUniqId(params UniqIdParams) string {
	var id string
	// Set prefix for unique id
	if params.Prefix != "" {
		id += params.Prefix
	}
	id += format(time.Now().Unix(), 8)
	// Increment global entropy value
	entropy++
	id += format(entropy, 5)
	// If we have more entropy add this
	if params.MoreEntropy == true {
		number := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * 10
		id += strconv.FormatFloat(number, 'E', -1, 64)[0:10]
	}

	return id
}

func format(number int64, width int) string {
	hex := strconv.FormatInt(number, 16)

	if width <= len(hex) {
		// so long we split
		return hex[0:width]
	}

	for len(hex) < width {
		hex = "0" + hex
	}

	return hex
}
