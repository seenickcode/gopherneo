package gopherneo

import (
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}
}
