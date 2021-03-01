package strutil

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	mathRand "math/rand"
	"time"
)

const (
	AlphaNumber = "0123456789"
	AlphaBet    = "abcdefghijklmnopqrstuvwxyz"
	AlphaNum    = "abcdefghijklmnopqrstuvwxyz0123456789"
	AlphaNum2   = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func Sha1(src interface{}) string {
	h := sha1.New()
	if s, ok := src.(string); ok {
		h.Write([]byte(s))
	} else {
		h.Write([]byte(fmt.Sprint(src)))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(src interface{}) string {
	h := sha256.New()
	if s, ok := src.(string); ok {
		h.Write([]byte(s))
	} else {
		h.Write([]byte(fmt.Sprint(src)))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Md5 Generate a 32-bit md5 string
func Md5(src interface{}) string {
	return GenMd5(src)
}

func Md5File(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	h := md5.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

// GenMd5 Generate a 32-bit md5 string
func GenMd5(src interface{}) string {
	h := md5.New()
	if s, ok := src.(string); ok {
		h.Write([]byte(s))
	} else {
		h.Write([]byte(fmt.Sprint(src)))
	}

	return hex.EncodeToString(h.Sum(nil))
}

func RandomBetween(min, max int64) string {
	var result int64
	if min > max || min == 0 || max == 0 {
		result = max
	}
	result = mathRand.Int63n(max-min) + min
	return fmt.Sprintf("%d", result)
}

func RandomNumbers(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(9) // 0 - 25
		cs[i] = AlphaNumber[idx]
	}

	return string(cs)
}

// RandomChars generate give length random chars at `a-z`
func RandomChars(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(25) // 0 - 25
		cs[i] = AlphaBet[idx]
	}

	return string(cs)
}

// RandomCharsV2 generate give length random chars in `0-9a-z`
func RandomCharsV2(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(35) // 0 - 35
		cs[i] = AlphaNum[idx]
	}

	return string(cs)
}

// RandomCharsV3 generate give length random chars in `0-9a-zA-Z`
func RandomCharsV3(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(61) // 0 - 61
		cs[i] = AlphaNum2[idx]
	}

	return string(cs)
}

// RandomBytes generate
func RandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandomString generate.
// Example:
// 	// this will give us a 44 byte, base64 encoded output
// 	token, err := RandomString(32)
// 	if err != nil {
//     // Serve an appropriately vague error to the
//     // user, but log the details internally.
// 	}
func RandomString(length int) (string, error) {
	b, err := RandomBytes(length)
	return base64.URLEncoding.EncodeToString(b), err
}
