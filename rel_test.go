package gopherneo

import (
	"testing"
	"time"
)

func TestLinkUnlinkNodes(t *testing.T) {

	db, err := NewConnectionWithToken("http://localhost:7474", "cee35b356a500f6bfd640146b4f3a771")

	db.DeleteNodes("Thing", "", "")

	// create node 1
	name1 := "joebob1"
	props1 := &map[string]interface{}{
		"name": name1,
	}
	err = db.CreateNode("Thing", props1, nil)
	if err != nil {
		t.Error(err)
	}
	// create node 2
	name2 := "joebob2"
	props2 := &map[string]interface{}{
		"name": name2,
	}
	err = db.CreateNode("ThingOther", props2, nil)
	if err != nil {
		t.Error(err)
	}

	// clear possible existing links
	err = db.UnlinkAllNodes("Thing", "name", name2, "LINKS_TO", "ThingOther")
	if err != nil {
		t.Error(err)
	}

	// link nodes
	timestamp1 := float64(time.Now().UnixNano()) // convert int64 -> float64
	relProps := &map[string]interface{}{
		"timestamp": timestamp1,
	}
	thingRel := &ThingLinksToThingRel{}
	err = db.LinkNodes("Thing", "name", name1, "ThingOther", "name", name2, "LINKS_TO", relProps, &thingRel)
	if err != nil {
		t.Error(err)
	}
	if thingRel.Timestamp != timestamp1 {
		t.Errorf("timestamp for rel doesn't match, was '%v', should be '%v'", thingRel.Timestamp, timestamp1)
	}

	// ensure they're linked
	linked, err := db.FindAllRelNodesPaginated("Thing", "name", name1, "ThingOther", "LINKS_TO", true, "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(linked.Rows) == 0 {
		t.Errorf("expected 1 linked node, got: %v", linked.Rows)
	}

	// link nodes without rel props
	err = db.LinkNodes("Thing", "name", name1, "ThingOther", "name", name2, "LINKS_TO", nil, nil)
	if err != nil {
		t.Error(err)
	}

	// ensure they're linked
	linked, err = db.FindAllRelNodesPaginated("Thing", "name", name1, "ThingOther", "LINKS_TO", true, "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(linked.Rows) == 0 {
		t.Errorf("expected 2 linked nodes, got: %v", linked.Rows)
	}

	// delink nodes
	err = db.UnlinkAllNodes("Thing", "name", name1, "LINKS_TO", "ThingOther")
	if err != nil {
		t.Error(err)
	}

	// ensure nothing else is linked now
	linked, err = db.FindAllRelNodesPaginated("Thing", "name", name1, "ThingOther", "LINKS_TO", true, "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(linked.Rows) != 0 {
		t.Errorf("exptected 0 nodes to be linked, got: %v", linked.Rows)
	}

	// cleanup
	err = db.DeleteNodes("Thing", "name", name1)
	if err != nil {
		t.Error(err)
	}
}
