package internal

import (
	"testing"
)

func TestIsDateValid(t *testing.T) {
	if IsValidDate("2023/03/05") != true {
		t.Error("expected date to be valid")
	}

	if IsValidDate("05/03/2023") != true {
		t.Error("expected date to be valid")
	}
}
