package ddos

import (
	"context"
	"faraway/wow/pkg/test"
	"testing"
	"time"
)

func Test_Reset_WithEmptyWindow(t *testing.T) {
	t.Parallel()

	// Initialize DDoSGuard with an empty window
	guard := NewGuard(0 * time.Second)

	// Set the current value to 10
	guard.current.Store(10)

	// Call reset and check the changes
	guard.reset()

	expected := int64(1)
	got, err := guard.IncRate(context.Background())

	test.Nil(t, "Inc error", err)

	test.Check(t, "Inc & reset with empty window", expected, got)
}

func Test_Reset_WithFullWindow(t *testing.T) {
	t.Parallel()

	// Initialize DDoSGuard with a 3-second window
	guard := NewGuard(3 * time.Second)

	// Fill the window completely
	for i := 1; i <= 3; i++ {
		guard.current.Store(int64(i * 10))
		guard.reset()
	}

	// Current value set to 40
	guard.current.Store(40)
	guard.reset()

	// Check the changes after window is full
	expected := int64(30 + 40 + 1) // last three values
	got, err := guard.IncRate(context.Background())

	test.Nil(t, "Inc error", err)
	test.Check(t, "Inc & reset with a full window, inc", expected, got)
}
