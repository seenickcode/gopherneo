package gopherneo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	c, err := NewConnection("localhost", "7474", "cee35b356a500f6bfd640146b4f3a771")
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

	db, err := NewConnection("localhost", "7474", "cee35b356a500f6bfd640146b4f3a771")
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
	err = json.Unmarshal(*cr.Rows[0][0], &newThing)
	if err != nil {
		t.Error(err)
	}
	if newThing.Name != nameVal {
		t.Errorf("name incorrect, should be '%v', is '%v'", nameVal, newThing.Name)
	}
	if newThing.Age != ageVal {
		t.Errorf("age incorrect, should be '%v', is '%v'", ageVal, newThing.Age)
	}

	db.DeleteNodes("Thing", "", "")
}

func TestReturnMultiNodes(t *testing.T) {

	db, err := NewConnection("localhost", "7474", "cee35b356a500f6bfd640146b4f3a771")
	if err != nil {
		t.Error(err)
	}

	name1 := "name1"
	name2 := "name2"
	cypher := fmt.Sprintf(`
		CREATE (t1:Thing { name: '%v' })
		CREATE (t2:Thing { name: '%v' })
		RETURN t1, t2`, name1, name2)
	params := &map[string]interface{}{}
	cr, err := db.ExecuteCypher(cypher, params)
	if err != nil {
		t.Error(err)
	}
	if len(cr.Rows) != 1 {
		t.Errorf("expected one row returned, query was: %v", cypher)
	}
	if len(cr.Rows[0]) != 2 {
		t.Errorf("expected two nodes returned in first row, query was: %v", cypher)
	}

	t1 := &Thing{}
	err = json.Unmarshal(*cr.Rows[0][0], &t1)
	if err != nil {
		t.Error(err)
	}

	t2 := &Thing{}
	err = json.Unmarshal(*cr.Rows[0][1], &t2)
	if err != nil {
		t.Error(err)
	}

	if t1.Name != name1 {
		t.Errorf("name incorrect, should be '%v', is '%v'", name1, t1.Name)
	}
	if t2.Name != name2 {
		t.Errorf("name incorrect, should be '%v', is '%v'", name2, t2.Name)
	}

	db.DeleteNodes("Thing", "", "")
}

// helpers

func errIfBlank(t *testing.T, key string, val string) {
	if len(val) == 0 {
		t.Error(fmt.Errorf("'%v' was blank", key))
	}
	return
}
