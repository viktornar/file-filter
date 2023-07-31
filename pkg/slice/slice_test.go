package slice

import "testing"

func TestSliceRemove(t *testing.T) {
	slice := []string{"a", "b"}
	expectedSlice := []string{"b"}

	slice = Remove[string](slice, 0)

	for idx, el := range slice {
		if expectedSlice[idx] != el {
			t.Errorf("expected element %s to be %s", el, expectedSlice[idx])
		}
	}
}

func TestSliceIndexOf(t *testing.T) {
	slice := []string{"a", "b"}

	for idx, el := range slice {
		if IndexOf[string](slice, el) != idx {
			t.Errorf("expected element %s to be at index %d", el, idx)
		}
	}
}
