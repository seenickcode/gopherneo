package gopherneo

import (
	//	"bytes"
	// "encoding/json"
	// "net/http"
	// "net/url"
	"regexp"
)

type Node struct {
	Data    map[string]interface{} `json:"data"`
	SelfURI string                 `json:"self"`
	//LabelsURI     string                 `json:"labels"`
	//Extensions    map[string]interface{} `json:"extensions"`
	//PropertiesURI string                 `json:"properties"`
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

// // GetNode fetches a node by its ID
// // http://docs.neo4j.org/chunked/stable/rest-api-nodes.html#rest-api-get-node
// func (c *Connection) GetNode(id string) (n *Node, err error) {

// 	// construct request
// 	// TODO use int instead of string for id param
// 	uri := joinPath([]string{c.NodeURI, id})
// 	req, _ := http.NewRequest("GET", uri, nil)
// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-Type", "application/json")

// 	// make request
// 	data, err := c.PerformRequest(req)
// 	if err != nil {
// 		return
// 	}

// 	// convert response data to node
// 	err = json.Unmarshal(data, &n)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// // GetNodesByLabelAndProperty gets nodes by label and property via the official Neo4j endpoint
// // http://docs.neo4j.org/chunked/stable/rest-api-node-labels.html#rest-api-get-nodes-by-label-and-property
// func (c *Connection) GetNodesByLabelAndProperty(label string, key string, value string) (nodes []Node, err error) {

// 	// prepare request
// 	qstring := url.Values{}
// 	qstring.Set(key, "\""+value+"\"")
// 	uri := joinPath([]string{c.Uri, "label", label, "nodes?" + qstring.Encode()})
// 	req, err := http.NewRequest("GET", uri, nil)
// 	if err != nil {
// 		return
// 	}
// 	req.Header.Add("Accept", "application/json; charset=UTF-8")

// 	// make request
// 	data, err := c.PerformRequest(req)
// 	if err != nil {
// 		return
// 	}

// 	// convert response to []Node
// 	err = json.Unmarshal(data, &nodes)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// // // CreateNode creates a node with a label and list of properties.
// // // Neo4j docs: http://docs.neo4j.org/chunked/stable/rest-api-nodes.html#rest-api-get-node
// // func (c *Connection) CreateNode(props interface{}) (n *Node, err error) {

// // 	// TODO: Verify that there's no way to specify label(s), as it seems
// // 	// we are forced to use Cypher to create a node with a label.

// // 	// convert properties to []byte
// // 	reqData, err := json.Marshal(props)
// // 	if err != nil {
// // 		return
// // 	}
// // 	reqBuf := bytes.NewReader(reqData)

// // 	// construct request
// // 	uri := c.NodeURI
// // 	req, err := http.NewRequest("POST", uri, reqBuf)
// // 	if err != nil {
// // 		return
// // 	}
// // 	req.Header.Add("Accept", "application/json; charset=UTF-8")
// // 	req.Header.Add("Content-Type", "application/json")

// // 	// make request
// // 	data, err := c.PerformRequest(req) // TODO ensure response is a 201?
// // 	if err != nil {
// // 		return
// // 	}

// // 	// convert response body to node
// // 	err = json.Unmarshal(data, &n)
// // 	if err != nil {
// // 		return
// // 	}
// // 	return
// // }

// // CreateNode creates a node with a label and list of properties.
// func (c *Connection) CreateNode(label string, props interface{}, unique bool) (n *Node, err error) {

// 	nodeAlias := "n"
// 	parts := make([]string, 3)

// 	// construct cypher
// 	if unique {
// 		parts[0] = "MERGE"
// 	} else {
// 		parts[0] = "CREATE"
// 	}

// 	// ** now that we should and are taking in an interface{},
// 	// pass props as a parameterized cypher query.
// 	// to do that, we'll need to change the cypher query
// 	// and include the props with the POST request

// 	//parts[1] = cypherizeNode(label, props, nodeAlias) // "(n:LABEL { PROPS })"
// 	parts[2] = "RETURN " + nodeAlias

// 	// perform query
// 	query := joinUsing(parts, " ")
// 	neoRes, err := c.Query(query) // TODO **use query with params**
// 	if err != nil {
// 		return
// 	}

// 	// convert response to node
// 	err = json.Unmarshal(*neoRes.Data[0][0], &n)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

func (n *Node) ID() string {

	// TODO return an int here
	re := regexp.MustCompile("db/data/node/(.+$)")
	matches := re.FindStringSubmatch(n.SelfURI)
	if matches != nil && len(matches) > 1 {
		return matches[1]
	}
	return ""
}
