package gopherneo

import "testing"

func TestQueryForCreate(t *testing.T) {

	db, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	// construct query which creates a thing
	// and returns a list of fields
	cypher1 := `
		CREATE (t:Thing { myprops }) 
		RETURN id(t) as id, t.name as name, t.age as age`

	props1 := struct {
		Name string `json:"name"`
		Age  int    `json:"age,int"`
	}{}
	props1.Name = "4379473927489327424343"
	props1.Age = 46

	query1 := NewQuery(cypher1)
	query1.Params["myprops"] = props1

	// perform query and print out list of values
	rows, err := db.Query(query1)
	if err != nil {
		t.Error(err)
	}

	if len(rows) != 1 {
		t.Errorf("returned rows not 1, query was: %v", cypher1)
	}

	// TODO test
	//log.Printf("> %v %T\n", row, row)

	// // ensure values are accurate
	// if row["id"].(float64) < 0 ||
	// 	row["name"].(string) != props1.Name {
	// 	//row["age"].(int) != props1.Age // TODO should be returned as int
	// 	t.Errorf("row data invalid, got: %v\n", row)
	// }
	// log.Printf("cool, we created a Thing node and got some fields back\n")
}
