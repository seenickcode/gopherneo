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

// 
// UTILITIES
//

func assertOk(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}