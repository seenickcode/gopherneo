package gopherneo

import (
	"regexp"
	"strconv"
)

type Node struct {
	Data            map[string]interface{} `json:"data"`
	SelfURI         string                 `json:"self"`
	LabelsURI       string                 `json:"labels"`
	PropertiesURI   string                 `json:"properties"`
	OutgoingRelsURI string                 `json:"outgoing_relationships"`
	IncomingRelsURI string                 `json:"incoming_relationships"`
	AllRelsURI      string                 `json:"all_relationships"`
}

// ID determines the Neo4j ID of a node.
// Be warned that Neo4j can change an ID at any time.
// Do not rely on this value.
func (n *Node) ID() int {
	re := regexp.MustCompile("db/data/node/(.+$)")
	matches := re.FindStringSubmatch(n.SelfURI)
	if matches != nil && len(matches) > 1 {
		id, err := strconv.Atoi(matches[1])
		if err == nil {
			return id
		}
	}
	return -1
}
