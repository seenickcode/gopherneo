package gopherneo

import (
	"fmt"
	"testing"
)

func TestNodes(t *testing.T) {

	neo, err := NewConnection("http://localhost:7474/db/data")
	if err != nil {
		t.Error(err)
	}

	// get node by ID
	nodeID := 0
	fmt.Println("fetching node %v", nodeID)
	_, err = neo.GetNode(nodeID)
	//fmt.Printf("node: %v %v\n", node, err)
	if err != nil {
		t.Error(err)
	}

	// get nodes by label and properties
	// TODO convert to generic values and add create from this
	nodes, err := neo.GetNodesByLabelAndProperty("User", "username", "nickTribeca")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf(">>> nodes: %v\n", nodes)

}
