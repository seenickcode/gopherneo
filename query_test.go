package gopherneo

import (
	"log"
	"testing"
)

func TestQueryForCreate(t *testing.T) {

	log.Println("testing our Queries")

	neo, err := NewConnection("http://localhost:7474/db/data")
	assertOk(t, err)

	// construct query which creates a thing
	// and returns a list of fields
	cypher1 := `
		CREATE (t:Thing { thingprops1 }) 
		RETURN id(t) as id, t.name as name, t.age as age`

	props1 := make(map[string]interface{})
	props1["name"] = "437947392748932742"
	props1["age"] = "46"
	query1 := NewQuery(cypher1, props1)
	query1.Params["thingprops1"] = props1

	// perform query and print out list of values
	resp, err := neo.Query(query1)
	assertOk(t, err)
	for _, result := range resp.Results {
		if result.ColNames[0] != "id" ||
			result.ColNames[1] != "name" ||
			result.ColNames[2] != "age" {
			t.Errorf("invalid col names: %v\n", result.ColNames)
		}
		if result.Rows[0].Data[1] != props1["name"] ||
			result.Rows[0].Data[2] != props1["age"] {
			t.Errorf("invalid col values: %v\n", result.Rows[0].Data)
		}
	}
	log.Printf("cool, we created a Thing node and got some fields back\n")

	log.Println()
}
