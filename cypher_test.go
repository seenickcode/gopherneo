package gopherneo

import (
	"fmt"
	"testing"
)

func TestPerformQuery(t *testing.T) {

	// create connection
	neo, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	// perform a cypher query
	query1 := "CREATE (t:Thing { prop1: 'val' }) RETURN t"
	fmt.Printf("performing query: %v\n", query1)
	cypResp, err := neo.PerformQuery(query1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("cypResp: %v\n", cypResp.Data)

}
