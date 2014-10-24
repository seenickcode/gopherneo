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
	props1["age"] = 46.0
	query1 := NewQuery(cypher1, props1)
	query1.Params["thingprops1"] = props1

	// perform query and print out list of values
	rows, err := neo.Query(query1)
	assertOk(t, err)

	row := rows[0]
	log.Printf("> %v %T\n", row, row)

	// ensure values are accurate
	if row["id"].(float64) < 0 ||
		row["name"].(string) != props1["name"] ||
		row["age"].(float64) != props1["age"] {
		t.Errorf("row data invalid, got: %v\n", row)
	}
	log.Printf("cool, we created a Thing node and got some fields back\n")

	log.Println()
}
