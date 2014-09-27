package gopherneo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Connection struct {
	Uri            string
	NodeURI        string      `json:"node"`
	RefNodeURI     string      `json:"reference_node"`
	NodeIndexURI   string      `json:"node_index"`
	RelIndexURI    string      `json:"relationship_index"`
	ExtInfoURI     string      `json:"extensions_info"`
	RelTypesURI    string      `json:"relationship_types"`
	BatchURI       string      `json:"batch"`
	CypherURI      string      `json:"cypher"`
	TransactionURI string      `json:"transaction"`
	VersionInfo    string      `json:"neo4j_version"`
	ExtensionsInfo interface{} `json:"extensions"`
}

func NewConnection(uri string) (conn *Connection, err error) {

	conn = &Connection{}

	// prepare request
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// perform request
	res, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Printf("> %v\n\n", res.Body)
	// get bytes from body
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	// deserialize
	err = json.Unmarshal(data, &conn)
	if err != nil {
		return
	}
	return
}
