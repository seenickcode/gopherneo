package gopherneo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"bytes"
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

// CreateNode creates a node with a list of properties
// Note that this function may be removed as the Neo4j is unrealistic. No one ever
// creates a new node without a label.
// http://docs.neo4j.org/chunked/stable/rest-api-nodes.html#rest-api-get-node
func (c *Connection) CreateNode(props interface{}) (n *Node, err error) {

	// compose URI
	uri := c.NodeURI

	log.Printf("creating node via: %v\n", uri)

	// prepare request body
	reqData, err := json.Marshal(props)
log.Printf(">>>>>>>>>>>> %v", reqData)	
	if err != nil {
  	return
  }
	reqBuf := bytes.NewReader(reqData)
log.Printf(">>>>>>>>>>>> %v", reqData)	

	// prepare request
	req, err := http.NewRequest("POST", uri, reqBuf)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json; charset=UTF-8")
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

	// TODO ensure response is a 201

	// deserialize
	err = json.Unmarshal(data, &n)
	if err != nil {
		return
	}

	// example response
	// {
	//   "extensions" : {
	//   },
	//   "outgoing_relationships" : "http://localhost:7474/db/data/node/287/relationships/out",
	//   "labels" : "http://localhost:7474/db/data/node/287/labels",
	//   "all_typed_relationships" : "http://localhost:7474/db/data/node/287/relationships/all/{-list|&|types}",
	//   "traverse" : "http://localhost:7474/db/data/node/287/traverse/{returnType}",
	//   "self" : "http://localhost:7474/db/data/node/287",
	//   "property" : "http://localhost:7474/db/data/node/287/properties/{key}",
	//   "properties" : "http://localhost:7474/db/data/node/287/properties",
	//   "outgoing_typed_relationships" : "http://localhost:7474/db/data/node/287/relationships/out/{-list|&|types}",
	//   "incoming_relationships" : "http://localhost:7474/db/data/node/287/relationships/in",
	//   "create_relationship" : "http://localhost:7474/db/data/node/287/relationships",
	//   "paged_traverse" : "http://localhost:7474/db/data/node/287/paged/traverse/{returnType}{?pageSize,leaseTime}",
	//   "all_relationships" : "http://localhost:7474/db/data/node/287/relationships/all",
	//   "incoming_typed_relationships" : "http://localhost:7474/db/data/node/287/relationships/in/{-list|&|types}",
	//   "metadata" : {
	//     "id" : 287,
	//     "labels" : [ ]
	//   },
	//   "data" : {
	//     "foo" : "bar"
	//   }
	// }

	return
}
