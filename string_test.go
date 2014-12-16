package gopherneo

import (
	"testing"
)

func TestStringJoining(t *testing.T) {

	// join a string as a path
	joined := joinPath([]string{"resources", "123"})
	if joined != "resources/123" {
		t.Error("joined result was: %v", joined)
	}

	// join a string
	joined = join([]string{"keanu", " ", "reeves"})
	if joined != "keanu reeves" {
		t.Error("joined result was: %v", joined)
	}

	// join with a delimiter
	joined = joinUsing([]string{"this", "that"}, "#")
	if joined != "this#that" {
		t.Error("joined result was: %v", joined)
	}
}

// func TestCypherGeneration(t *testing.T) {

// 	// create cypher node from parameters
// 	props := make(map[string]interface{})
// 	props["car"] = "My Car"
// 	props["bar"] = 3
// 	//props["jar"] = 45.0 // TODO support floats
// 	cypherizedNode := propsToCypherString("Thing", props, "t")
// 	if cypherizedNode != "(t:Thing { car: 'My Car', bar: 3 })" {
// 		t.Error("cypherized node is inaccurate: %v\n", cypherizedNode)
// 	}
// }
