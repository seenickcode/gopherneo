package gopherneo

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	neo, err := NewConnection("http://localhost:7474/db/data")
	fmt.Printf("%v %v", neo, err)
	if err != nil {
		t.Error(err)
	}
}
