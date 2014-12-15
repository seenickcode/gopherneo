package gopherneo

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	c, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}
	errIfBlank(t, "Uri", c.Uri)
	errIfBlank(t, "NodeURI", c.NodeURI)
	errIfBlank(t, "NodeLabelsURI", c.NodeLabelsURI)
	errIfBlank(t, "CypherURI", c.CypherURI)
	errIfBlank(t, "TransactionURI", c.TransactionURI)
}

// helpers

func errIfBlank(t *testing.T, key string, val string) {
	if len(val) == 0 {
		t.Error(fmt.Errorf("'%v' was blank", key))
	}
	return
}
