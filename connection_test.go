package gopherneo

import (
	"encoding/json"
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

func TestQueryWithProps(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	// construct query which creates a thing and returns the node

	type Thing struct {
		Name string `json:"name"`
		Age  int    `json:"age,int"`
	}

	// prepare the cypher statement
	cypher1 := `
		CREATE (t:Thing { myprops }) 
		RETURN t`

	// prepare the cypher props for "myprops"
	nameVal := "4379473927489327424343"
	ageVal := 46
	props := &Thing{Name: nameVal, Age: ageVal}

	// add my cypher props to a map[string]interface{}
	params := &map[string]interface{}{
		"myprops": props,
	}

	cr, err := db.ExecuteCypher(cypher1, params)
	if err != nil {
		t.Error(err)
	}
	if len(cr.Rows) != 1 {
		t.Errorf("returned rows not 1, query was: %v", cypher1)
	}

	newThing := &Thing{}
	err = json.Unmarshal(cr.Rows[0], &newThing)
	if err != nil {
		t.Error(err)
	}
	if newThing.Name != nameVal {
		t.Errorf("name incorrect, should be '%v', is '%v'", nameVal, newThing.Name)
	}
	if newThing.Age != ageVal {
		t.Errorf("age incorrect, should be '%v', is '%v'", ageVal, newThing.Age)
	}
}

// helpers

func errIfBlank(t *testing.T, key string, val string) {
	if len(val) == 0 {
		t.Error(fmt.Errorf("'%v' was blank", key))
	}
	return
}
