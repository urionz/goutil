package goutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urionz/goutil"
)

func TestFuncName(t *testing.T) {
	name := goutil.FuncName(goutil.PkgName)
	assert.Equal(t, "github.com/urionz/goutil.PkgName", name)
}

func TestPkgName(t *testing.T) {
	name := goutil.PkgName()
	assert.Equal(t, "goutil", name)
}

func TestPanicIfErr(t *testing.T) {
	goutil.PanicIfErr(nil)
}
