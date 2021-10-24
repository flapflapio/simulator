package tm

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimulateUnimplemented(t *testing.T) {
	panicChan := make(chan string)
	go func() {
		defer func() {
			panicChan <- fmt.Sprintf("%v", recover())
		}()
		(&TM{}).Simulate("")
	}()

	select {
	case e := <-panicChan:
		assert.Equal(t, "not implemented", e)
	case <-time.After(1 * time.Second):
		t.Error("timed out (1s) while waiting for TM.Simulate to panic")
	}
}
