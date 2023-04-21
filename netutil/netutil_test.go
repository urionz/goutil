package netutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urionz/goutil/netutil"
)

func TestInternalIP(t *testing.T) {
	assert.NotEmpty(t, netutil.InternalIP())
}
