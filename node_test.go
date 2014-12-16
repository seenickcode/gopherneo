package gopherneo

import "testing"

func TestCreateNodeWithLabel(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	type Thing struct {
		Name string `json:"name"`
		Age  int    `json:age,int`
	}

	name1 := "joebob2"
	age1 := 49
	props := &map[string]interface{}{
		"name": name1,
		"age":  age1,
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
}
