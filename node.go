package gopherneo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Node struct {
	Labels        string                 `json:"labels"`
	Data          map[string]interface{} `json:"data"`
	Extensions    map[string]interface{} `json:"extensions"`
	PropertiesURI string                 `json:"properties"`
	// OutgoingRelsURI      string                 `json:"outgoing_relationships"`
	// TraverseURI          string                 `json:"traverse"`
	// AllTypedRelsURI      string                 `json:"all_typed_relationships"`
	// OutgoingURI          string                 `json:"outgoing_typed_relationships"`
	// IncomingRelsURI      string                 `json:"incoming_relationships"`
	// CreateRelURI         string                 `json:"create_relationship"`
	// PagedTraverseURI     string                 `json:"paged_traverse"`
	// AllRelsURI           string                 `json:"all_relationships"`
	// IncomingTypedRelsURI string                 `json:"incoming_typed_relationships"`
}

// GetNode fetches a node by its ID
// http://docs.neo4j.org/chunked/stable/rest-api-nodes.html#rest-api-get-node
func (c *Connection) GetNode(id int) (n *Node, err error) {

	// compose URI
	uri := joinPath([]string{c.NodeURI, strconv.Itoa(id)})

	log.Printf("fetching node via: %v\n", uri)

	// prepare request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// perform request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	// get bytes from body
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}

	// deserialize
	err = json.Unmarshal(data, &n)
	if err != nil {
		return
	}
	return
}

// GetNodesByLabelAndProperty gets nodes by label and property via the official Neo4j endpoint
// http://docs.neo4j.org/chunked/stable/rest-api-node-labels.html#rest-api-get-nodes-by-label-and-property
func (c *Connection) GetNodesByLabelAndProperty(label string, key string, value string) (nodes []Node, err error) {

	// compose URI
	// FIXME use a proper string join
	qstring := url.Values{}
	qstring.Set(key, "\""+value+"\"")
	uri := joinPath([]string{c.Uri, "label", label, "nodes?" + qstring.Encode()})

	log.Printf("fetching node by label and properties via: %v\n", uri)

	// prepare request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json; charset=UTF-8")

	// perform request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	// get bytes from body
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}

	// deserialize
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return
	}
	return
}
