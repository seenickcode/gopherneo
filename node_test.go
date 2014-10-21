package gopherneo

import (
	"log"

	"strconv"
	"testing"
	"time"
)

type TestUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

func TestNodes(t *testing.T) {

	log.Println("testing our Node queries")

	neo, err := NewConnection("http://localhost:7474/db/data")
	assertOk(t, err)

	timestamp := strconv.Itoa(time.Now().Nanosecond())

	// create a new User via CreateNode
	// user1Props := make(map[string]interface{})
	// user1Props["username"] = "testusername" + timestamp
	// user1Props["name"] = "Tester1" + timestamp
	// user1Props["age"] = 45
	user1 := &TestUser{
		Username: "testusername" + timestamp, 
		Name: "Tester1" + timestamp,
		Age: 45
	}
	log.Printf("creating a User node with username: %v\n", user1Props)
	user1Node, err := neo.CreateNode("Thing", user1, true)

	// assert username
	if user1Node.Data["username"].(string) != user1Props["username"] {
		t.Error("create response 'username' incorrect")
	}
	if user1Node.Data["name"].(string) != user1Props["name"] {
		t.Error("create response 'name' incorrect")
	}
	// TODO **figure out why response is a float64 and not an int. golang is treating it as float64**
	// if user1Node.Data["age"].(float64).Convert(int) != user1Props["age"].(int) {
	// 	t.Error("create response 'age' incorrect")
	// }
	if user1Node.Data["age"].(float64) != 45 { // NOTE: **temporary**
		t.Error("create response 'age' incorrect")
	}
	// ensure we can determine a Node's ID
	if len(user1Node.ID()) == 0 {
		t.Error("couldn't determine Node's ID")
	}
	log.Printf("ok, we created User node %v\n", user1Node.ID())

	// get node by ID
	user1Node, err = neo.GetNode(user1Node.ID())
	log.Printf("fetching User by ID %v\n", user1Node.ID())
	assertOk(t, err)
	if len(user1Node.Data["username"].(string)) == 0 {
		t.Error("couldn't fetch User by ID")
	}
	log.Printf("ok, we got our Node: %v\n", user1Node)

	// // get nodes by label and properties
	// // TODO convert to generic values and add create from this
	// _, err = neo.GetNodesByLabelAndProperty("User", "username", "nickTribeca")
	// log.Printf("fetched node by label and property\n")
	// if err != nil {
	// 	t.Error(err)
	// }

	// clean up
	// TODO

	log.Println()
}
