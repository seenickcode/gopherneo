package gopherneo

import (
	"log"
	"testing"
)

func TestNodes(t *testing.T) {

	neo, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	// TODO create test nodes, then tear them down

	// create a new node
	type TestUser struct {
		Username string `json:"username"`
		Name string `json:"name"`
	}
	
	user1Props := &TestUser{Username: "testusername1", Name: "Tester1"}
	user1Node, err := neo.CreateNode(user1Props)
	log.Printf("created node with props: %v", user1Props)
	log.Printf("created node with props: %v", user1Node)

	// get node by ID
	nodeID := 0
	_, err = neo.GetNode(nodeID)
	log.Printf("fetched node: %v\n", nodeID)
	if err != nil {
		t.Error(err)
	}

	// get nodes by label and properties
	// TODO convert to generic values and add create from this
	_, err = neo.GetNodesByLabelAndProperty("User", "username", "nickTribeca")
	log.Printf("fetched node by label and property\n")
	if err != nil {
		t.Error(err)
	}
}
