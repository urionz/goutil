package process_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urionz/goutil/sysutil/process"
)

func TestProcessExists(t *testing.T) {
	pid := os.Getpid()

	assert.True(t, process.Exists(pid))
}

func TestPID(t *testing.T) {
	assert.True(t, process.PID() > 0)
}
