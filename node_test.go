package gopherneo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Thing struct {
	Name      string  `json:"name"`
	Age       int     `json:"age",int`
	CreatedAt float64 `json:"created_at",float64`
}

type ThingLinksToThingRel struct {
	Timestamp float64 `json:"timestamp",float64`
}

func UnmarshalThings(rows [][]*json.RawMessage) []Thing {
	things := make([]Thing, 2)
	for i, row := range rows {
		thing := &Thing{}
		err := json.Unmarshal(*row[0], thing)
		if err == nil {
			things[i] = *thing
		}
	}
	return things
}

func TestFindNode(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	// create a node
	name1 := "joebob7"
	props := &map[string]interface{}{
		"name": name1,
	}
	newThing := &Thing{}
	err = db.CreateNode("Thing", props, newThing)
	if err != nil {
		t.Error(err)
	}

	// find the node
	fetchedThing := &Thing{}
	err = db.FindNode("Thing", "name", name1, &fetchedThing)
	if err != nil {
		t.Error(err)
	}
	if fetchedThing.Name != newThing.Name {
		t.Errorf("created thing named '%v' didn't match fetched thing named '%v'", newThing.Name, fetchedThing.Name)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name1)
	if err != nil {
		t.Error(err)
	}
}

func TestFindNodesWithValuesPaginated(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	// cleanup possibly preexisting nodes
	_ = db.DeleteNodes("Thing", "", "")

	numNodes := 5
	for i := 0; i < numNodes; i++ {
		// create node
		name := fmt.Sprintf("joebobby%d", i)
		props := &map[string]interface{}{
			"name": name,
		}
		err = db.CreateNode("Thing", props, nil)
		if err != nil {
			t.Error(err)
		}
	}

	// get all nodes
	rows, err := db.FindNodesWithValuePaginated("Thing", "", "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(rows) != numNodes {
		t.Errorf("found %d nodes, expected %d: %v", len(rows), numNodes, rows)
	}

	// ensure pagination page 1 results are accurate
	rows, err = db.FindNodesWithValuePaginated("Thing", "", "", 0, 2)
	things := UnmarshalThings(rows)
	expectedName1 := "joebobby1"
	if things[1].Name != expectedName1 {
		t.Errorf("expected paginated result to be called '%v' and it's called '%v'", expectedName1, things[1].Name)
	}
	// ensure pagination page 2 results are accurate
	rows, err = db.FindNodesWithValuePaginated("Thing", "", "", 1, 2)
	things = UnmarshalThings(rows)
	expectedName2 := "joebobby3"
	if things[1].Name != expectedName2 {
		t.Errorf("expected paginated result to be called '%v' and it's called '%v'", expectedName2, things[1].Name)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNode(t *testing.T) {

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

	err = db.CreateNode("Thing", props, newThing)
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

func TestLinkNodes(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	db.DeleteNodes("Thing", "", "")

	// create node 1
	name1 := "joebob1"
	props1 := &map[string]interface{}{
		"name": name1,
	}
	err = db.CreateNode("Thing", props1, nil)
	if err != nil {
		t.Error(err)
	}
	// create node 2
	name2 := "joebob2"
	props2 := &map[string]interface{}{
		"name": name2,
	}
	err = db.CreateNode("Thing", props2, nil)
	if err != nil {
		t.Error(err)
	}

	// link nodes
	timestamp1 := float64(time.Now().UnixNano()) // convert int64 -> float64
	relProps := &map[string]interface{}{
		"timestamp": timestamp1,
	}
	thingRel := &ThingLinksToThingRel{}
	err = db.LinkNodes("Thing", "name", name1, "Thing", "name", name2, "LINKS_TO", relProps, &thingRel)
	if err != nil {
		t.Error(err)
	}
	if thingRel.Timestamp != timestamp1 {
		t.Errorf("timestamp for rel doesn't match, was '%v', should be '%v'", thingRel.Timestamp, timestamp1)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name1)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNodeErrors(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	type Thing struct {
		Name      string  `json:"name"`
		Age       int     `json:"age",int`
		CreatedAt float64 `json:"created_at",float64`
	}

	emptyProps := &map[string]interface{}{}

	newThing := &Thing{}

	err = db.CreateNode("", emptyProps, newThing)
	if err == nil {
		t.Errorf("should have been a 'no label' error")
	}

	err = db.CreateNode("Thing", emptyProps, newThing)
	if err != nil {
		t.Error(err)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "", "")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateNode(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")

	name1 := "joebob5"
	age1 := 46
	name2 := "joebob6"
	props1 := &map[string]interface{}{
		"name": name1,
		"age":  age1,
	}

	newThing1 := &Thing{}

	// create a node
	err = db.CreateNode("Thing", props1, newThing1)
	if err != nil {
		t.Error(err)
	}

	props2 := &map[string]interface{}{
		"name": name2,
	}

	updatedThing1 := &Thing{}

	// update it
	err = db.UpdateNode("Thing", "name", name1, props2, updatedThing1)
	if err != nil {
		t.Error(err)
	}

	if updatedThing1.Name != name2 {
		t.Errorf("name doesn't match, was '%v', should be '%v'", updatedThing1.Name, name2)
	}
	if updatedThing1.Age != age1 {
		t.Errorf("age doesn't match, was '%v', should be '%v'", updatedThing1.Age, age1)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name2)
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
	err = db.CreateNode("Thing", props, newThing)
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
