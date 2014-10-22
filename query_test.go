package gopherneo

import (
	"log"
	"testing"
)

func TestQuery(t *testing.T) {

	log.Println("testing our Queries")

	neo, err := NewConnection("http://localhost:7474/db/data")
	assertOk(t, err)

	// construct query which creates a thing
	// and returns a list of fields
	cypher1 := "CREATE (t:Thing { props }) RETURN id(t) as id, t.name, t.age"
	props1 := make(map[string]interface{})
	props1["name"] = "437947392748932742"
	props1["age"] = "46"
	query1 := NewQuery(cypher1, props1)
	query1.Params["props"] = props1

	// result1 := struct {
	// 	ID   string `json:"id"`
	// 	Name string `json:"name"`
	// }{}

	// perform query
	resp, err := neo.Query(query1)
	assertOk(t, err)
	for _, result := range resp.Results {
		for _, row := range result.Rows {
			for colNdx, colData := range row.Data {
				colName := result.ColumnNames[colNdx]
				log.Printf(">>>>>>> %v: %v\n", colName, colData)
			}
		}
	}

	log.Println()
}
