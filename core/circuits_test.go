package core

import (
	"testing"
)

func TestExpectedStates(t *testing.T) {
	t.Run("Bell", func(t *testing.T) {
		got := ExpectedBellState()
		if len(got) != 4 {
			t.Errorf("expected 4 amplitudes, got %d", len(got))
		}
		// Sum of squared magnitudes should be 1
		sum := 0.0
		for _, c := range got {
			sum += real(c)*real(c) + imag(c)*imag(c)
		}
		if sum < 0.999 || sum > 1.001 {
			t.Errorf("expected total probability 1, got %f", sum)
		}
	})

	t.Run("GHZ", func(t *testing.T) {
		got := ExpectedGHZState(3)
		if len(got) != 8 {
			t.Errorf("expected 8 amplitudes, got %d", len(got))
		}
		if got[0] == 0 || got[7] == 0 {
			t.Errorf("expected non-zero amplitudes at 0 and 7")
		}
	})

	t.Run("QFT", func(t *testing.T) {
		got := ExpectedQFTState(2)
		if len(got) != 4 {
			t.Errorf("expected 4 amplitudes, got %d", len(got))
		}
		for i, c := range got {
			if real(c) == 0 {
				t.Errorf("amplitude %d should be non-zero", i)
			}
		}
	})
}
