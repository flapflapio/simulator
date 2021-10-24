package nfa

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests that `NFA.Simulate` is not implemented
func TestSimulateUnimplemented(t *testing.T) {
	panicChan := make(chan string)
	go func() {
		defer func() {
			panicChan <- fmt.Sprintf("%s", recover())
		}()
		(&NFA{}).Simulate("")
	}()

	select {
	case e := <-panicChan:
		assert.Equal(t, "not implemented", e)
	case <-time.After(1 * time.Second):
		t.Error("timed out (1s) while waiting for NFA.Simulate to panic")
	}
}
