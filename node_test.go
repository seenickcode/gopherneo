package gopherneo

import (
	"fmt"
	"testing"
	"time"
)

type Thing struct {
	Name      string  `json:"name"`
	Age       int     `json:"age",int`
	CreatedAt float64 `json:"created_at",float64`
}

func TestCreateNodeWithLabel(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	name1 := "joebob2"
	age1 := 49
	timestamp1 := time.Now().UnixNano()
	props := &map[string]interface{}{
		"name":       name1,
		"age":        age1,
		"created_at": timestamp1,
	}

	newThing := &Thing{}

	err = db.CreateNodeWithLabel("Thing", props, newThing)
	if err != nil {
		t.Error(err)
	}
	if newThing.Name != name1 {
		t.Errorf("name doesn't match, was '%v', should be '%v'", newThing.Name, name1)
	}
	if newThing.Age != age1 {
		t.Errorf("age doesn't match, was '%v', should be '%v'", newThing.Age, age1)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name1)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNodeWithLabelErrors(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	type Thing struct {
		Name      string  `json:"name"`
		Age       int     `json:"age",int`
		CreatedAt float64 `json:"created_at",float64`
	}

	emptyProps := &map[string]interface{}{}

	newThing := &Thing{}

	err = db.CreateNodeWithLabel("", emptyProps, newThing)
	if err == nil {
		t.Errorf("should have been a 'no label' error")
	}

	err = db.CreateNodeWithLabel("Thing", emptyProps, newThing)
	if err != nil {
		t.Error(err)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteNodes(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	// create node
	name1 := "joebob3"
	props := &map[string]interface{}{
		"name": name1,
	}

	newThing := &Thing{}
	err = db.CreateNodeWithLabel("Thing", props, newThing)
	if err != nil {
		t.Error(err)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name1)
	if err != nil {
		t.Error(err)
	}

	// TODO list nodes and ensure none are returned
}

func TestFindNodesWithValuesPaginated(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	numNodes := 5
	for i := 0; i < numNodes; i++ {
		// create node
		name := fmt.Sprintf("joebobby%d", i)
		props := &map[string]interface{}{
			"name": name,
		}
		err = db.CreateNodeWithLabel("Thing", props, nil)
		if err != nil {
			t.Error(err)
		}
	}

	rows, err := db.FindNodesWithValuePaginated("Thing", "", "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(rows) != 5 {
		t.Errorf("found %d nodes, expected %d: %v", len(rows), numNodes, rows)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "", "")
	if err != nil {
		t.Error(err)
	}
}
