package strutil

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"text/template"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
)

// Position for padding string
const (
	PosLeft uint8 = iota
	PosRight
)

// IsAlphabet char
func IsAlphabet(char uint8) bool {
	// A 65 -> Z 90
	if char >= 'A' && char <= 'Z' {
		return true
	}

	// a 97 -> z 122
	if char >= 'a' && char <= 'z' {
		return true
	}

	return false
}

/*************************************************************
 * String filtering
 *************************************************************/

// Trim string
func Trim(s string, cutSet ...string) string {
	if len(cutSet) > 0 && cutSet[0] != "" {
		return strings.Trim(s, cutSet[0])
	}

	return strings.TrimSpace(s)
}

// TrimLeft char in the string.
func TrimLeft(s string, cutSet ...string) string {
	if len(cutSet) > 0 {
		return strings.TrimLeft(s, cutSet[0])
	}

	return strings.TrimLeft(s, " ")
}

// TrimRight char in the string.
func TrimRight(s string, cutSet ...string) string {
	if len(cutSet) > 0 {
		return strings.TrimRight(s, cutSet[0])
	}

	return strings.TrimRight(s, " ")
}

// FilterEmail filter email, clear invalid chars.
func FilterEmail(s string) string {
	s = strings.TrimSpace(s)
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s
	}

	// According to rfc5321, "The local-part of a mailbox MUST BE treated as case sensitive"
	return s[0:i] + "@" + strings.ToLower(s[i+1:])
}

/*************************************************************
 * String operation
 *************************************************************/

// Split string to slice. will clear empty string node.
func Split(s, sep string) (ss []string) {
	if s = strings.TrimSpace(s); s == "" {
		return
	}

	for _, val := range strings.Split(s, sep) {
		if val = strings.TrimSpace(val); val != "" {
			ss = append(ss, val)
		}
	}
	return
}

// Substr for a string.
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	strLen := len(runes)

	// pos is to large
	if pos >= strLen {
		return ""
	}

	l := pos + length
	if l > strLen {
		l = strLen
	}

	return string(runes[pos:l])
}

// Padding a string.
func Padding(s, pad string, length int, pos uint8) string {
	diff := len(s) - length
	if diff >= 0 { // do not need padding.
		return s
	}

	if pad == "" || pad == " " {
		mark := ""
		if pos == PosRight { // to right
			mark = "-"
		}

		// padding left: "%7s", padding right: "%-7s"
		tpl := fmt.Sprintf("%s%d", mark, length)
		return fmt.Sprintf(`%`+tpl+`s`, s)
	}

	if pos == PosRight { // to right
		return s + Repeat(pad, -diff)
	}

	return Repeat(pad, -diff) + s
}

// PadLeft a string.
func PadLeft(s, pad string, length int) string {
	return Padding(s, pad, length, PosLeft)
}

// PadRight a string.
func PadRight(s, pad string, length int) string {
	return Padding(s, pad, length, PosRight)
}

// Repeat repeat a string
func Repeat(s string, times int) string {
	if times < 2 {
		return s
	}

	var ss []string
	for i := 0; i < times; i++ {
		ss = append(ss, s)
	}

	return strings.Join(ss, "")
}

// RepeatRune repeat a rune char.
func RepeatRune(char rune, times int) (chars []rune) {
	for i := 0; i < times; i++ {
		chars = append(chars, char)
	}
	return
}

// Replaces replace multi strings
//
// 	pairs: {old1: new1, old2: new2, ...}
//
// Can also use:
// 	strings.NewReplacer("old1", "new1", "old2", "new2").Replace(str)
func Replaces(str string, pairs map[string]string) string {
	ss := make([]string, len(pairs)*2)
	for old, newVal := range pairs {
		ss = append(ss, old, newVal)
	}

	return strings.NewReplacer(ss...).Replace(str)
}

// PrettyJSON get pretty Json string
// Deprecated
//  please use fmtutil.PrettyJSON() or jsonutil.Pretty() instead it
func PrettyJSON(v interface{}) (string, error) {
	out, err := json.MarshalIndent(v, "", "    ")
	return string(out), err
}

// RenderTemplate render text template
func RenderTemplate(input string, data interface{}, fns template.FuncMap, isFile ...bool) string {
	return RenderText(input, data, fns, isFile...)
}

// RenderText render text template
func RenderText(input string, data interface{}, fns template.FuncMap, isFile ...bool) string {
	t := template.New("simple-text")
	t.Funcs(template.FuncMap{
		// don't escape content
		"raw": func(s string) string {
			return s
		},
		"trim": func(s string) string {
			return strings.TrimSpace(string(s))
		},
		// join strings
		"join": func(ss []string, sep string) string {
			return strings.Join(ss, sep)
		},
		// lower first char
		"lcFirst": func(s string) string {
			return LowerFirst(s)
		},
		// upper first char
		"upFirst": func(s string) string {
			return UpperFirst(s)
		},
	})

	// custom add template functions
	if len(fns) > 0 {
		t.Funcs(fns)
	}

	if len(isFile) > 0 && isFile[0] {
		template.Must(t.ParseFiles(input))
	} else {
		template.Must(t.Parse(input))
	}

	// use buffer receive rendered content
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}

// UUID
func UUID() string {
	u, _ := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

var objectIdCounter uint32 = 0

var machineId = readMachineId()

type ObjectId string

func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

// GUID returns a new unique ObjectId.
// 4byte 时间，
// 3byte 机器ID
// 2byte pid
// 3byte 自增ID
func GUID() ObjectId {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectId(b[:])
}

// Hex returns a hex representation of the ObjectId.
// 返回16进制对应的字符串
func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}
