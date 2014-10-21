package gopherneo

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {

	log.Println("testing our Cypher queries")

	// create connection
	neo, err := NewConnection("http://localhost:7474/db/data")
	assertOk(t, err)

	// set up our queries
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	thingName := "name" + timestamp // TODO use epoch

	// create a unique Thing, return a single node
	query1 := "MERGE (t:Thing { name: '" + thingName + "' }) RETURN t"
	log.Printf("creating a Thing node\n")
	neoRes, err := neo.Query(query1)
	assertOk(t, err)

	// was our node created?
	thingNode := &Node{}
	err = json.Unmarshal(*neoRes.Data[0][0], &thingNode)
	assertOk(t, err)
	if thingNode.Data["name"] != thingName {
		t.Error(t, "'thingName' invalid, was: %v\n", thingNode.Data["name"])
	}
	log.Printf("nice, Neo created our Thing and responded with a Thing node\n")

	// clean up
	// TODO

	log.Println()
}
