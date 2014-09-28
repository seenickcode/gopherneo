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
}
