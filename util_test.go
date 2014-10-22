package gopherneo

import (
	"testing"
)

func TestJoin(t *testing.T) {

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

	// create a timestamp (epoch in nanoseconds)
	ts := generateTimestamp()
	if ts <= 0 {
		t.Error("couldn't create timestamp, was: %v\n", ts)
	}

	// create cypher node from parameters
	props := make(map[string]interface{})
	props["car"] = "My Car"
	props["bar"] = 3
	//props["jar"] = 45.0 // TODO support floats
	cypherizedNode := cypherizeNode("Thing", props, "t")
	if cypherizedNode != "(t:Thing { car: 'My Car', bar: 3 })" {
		t.Error("cypherized node is inaccurate: %v\n", cypherizedNode)
	}

}
