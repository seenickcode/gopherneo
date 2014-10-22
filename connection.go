package gopherneo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Connection struct {
	httpClient     *http.Client
	Uri            string
	Version        string `json:"neo4j_version"`
	NodeURI        string `json:"node"`
	NodeLabelsURI  string `json:"node_labels"`
	CypherURI      string `json:"cypher"`
	TransactionURI string `json:"transaction"`
	//  Extensions     interface{} `json:"extensions"`
	// RefNodeURI     string      `json:"reference_node"`
	// NodeIndexURI   string      `json:"node_index"`
	// RelIndexURI    string      `json:"relationship_index"`
	// ExtInfoURI     string      `json:"extensions_info"`
	// RelTypesURI    string      `json:"relationship_types"`
	// BatchURI       string      `json:"batch"`
}

// get the Neo4j "service root"
// http://docs.neo4j.org/chunked/stable/rest-api-service-root.html
func NewConnection(uri string) (c *Connection, err error) {

	c = &Connection{httpClient: &http.Client{}, Uri: uri}

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
	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}
	return
}

func (c *Connection) performRequest(req *http.Request) (data []byte, err error) {

	// perform request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	// get bytes from body
	data, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	return
}
