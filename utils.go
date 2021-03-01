package goutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/urionz/goutil/jsonutil"
	"golang.org/x/crypto/bcrypt"
)

// Go is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
// from beego/bee
func Go(f func() error) error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return <-ch
}

func EncodePassword(rawPassword string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return ""
	}
	return string(hash)
}

func ValidatePassword(encodePassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(inputPassword))
	return err == nil
}

// Filling filling a model from submitted data
// form 提交过来的数据结构体
// model 定义表模型的数据结构体
// 相当于是在合并两个结构体(data 必须是 model 的子集)
func Filling(form interface{}, model interface{}) error {
	jsonBytes, _ := jsonutil.Encode(form)
	return jsonutil.Decode(jsonBytes, model)
}

// FuncName get func name
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// PkgName get current package name
func PkgName() string {
	_, filePath, _, _ := runtime.Caller(0)
	file, _ := os.Open(filePath)
	r := bufio.NewReader(file)
	line, _, _ := r.ReadLine()
	pkgName := bytes.TrimPrefix(line, []byte("package "))

	return string(pkgName)
}

// PanicIfErr if error is not empty
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Contains(search interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == search {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(search)).IsValid() {
			return true
		}
	}
	return false
}

func ContainsIgnoreCase(search string, target []string) bool {
	if len(search) == 0 {
		return false
	}
	if len(target) == 0 {
		return false
	}
	search = strings.ToLower(search)
	for i := 0; i < len(target); i++ {
		if strings.ToLower(target[i]) == search {
			return true
		}
	}
	return false
}

func StructToMap(obj interface{}, excludes ...string) map[string]interface{} {
	var data = make(map[string]interface{})
	keys := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)
	fillMap(data, keys, values, excludes...)
	return data
}

func fillMap(data map[string]interface{}, keys reflect.Type, values reflect.Value, excludes ...string) {
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	if keys.Kind() == reflect.Ptr {
		keys = keys.Elem()
	}

	for i := 0; i < keys.NumField(); i++ {
		keyField := keys.Field(i)
		valueField := values.Field(i)

		if keyField.Anonymous {
			fillMap(data, keyField.Type, valueField, excludes...)
		} else {
			if !ContainsIgnoreCase(keyField.Name, excludes) {
				jsonTag := keyField.Tag.Get("json")
				if len(jsonTag) > 0 {
					data[jsonTag] = valueField.Interface()
				} else {
					data[keyField.Name] = valueField.Interface()
				}
			}
		}
	}
}

func MapToStruct(obj interface{}, data map[string]interface{}) error {
	for k, v := range data {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj ", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type ")
	}
	structFieldValue.Set(val)
	return nil
}

// 获取struct字段
func StructFields(s interface{}) []reflect.StructField {
	t := StructTypeOf(s)
	if t.Kind() != reflect.Struct {
		return nil
	}

	var results []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		results = append(results, f)
		// if f.Anonymous {
		// 	fields := StructFields(f.Type)
		// 	results = append(results, fields...)
		// }
	}
	return results
}

// 获取struct name
func StructName(s interface{}) string {
	t := StructTypeOf(s)
	return t.Name()
}

func StructTypeOf(s interface{}) reflect.Type {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

type stop struct {
	error
}

func RetryFunc(attempt int, fn func() error, sleep ...time.Duration) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempt--; attempt > 0 {
			if len(sleep) > 0 {
				time.Sleep(sleep[0])
				return RetryFunc(attempt, fn, sleep...)
			}
		}
		return err
	}
	return nil
}
