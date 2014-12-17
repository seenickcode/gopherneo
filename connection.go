package gopherneo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tideland/goas/v3/logger"
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

type TransactionStatement struct {
	Cypher      string                 `json:"statement"`
	Params      map[string]interface{} `json:"parameters"`
	ResultTypes []string               `json:"resultDataContents"`
}

type TransactionResponse struct {
	Results []TransactionResult `json:"results"`
	Errors  []TransactionError  `json:"errors"`
}

type TransactionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// represents a transactional response result
type TransactionResult struct {
	Columns []string                        `json:"columns"`
	Data    []map[string][]*json.RawMessage `json:"data"`
}

func (r *TransactionResult) RowData() (rows [][]*json.RawMessage) {
	rows = make([][]*json.RawMessage, len(r.Data))
	for i, resultType := range r.Data {
		if val, ok := resultType["row"]; ok {
			rows[i] = val
		}
	}
	return
}

// type TransactionResultData struct {
// 	RowData   []*json.RawMessage `json:"row"`
// 	RestData  []*json.RawMessage `json:"rest"`
// 	GraphData []*json.RawMessage `json:"graph"`
// }

// get the Neo4j "service root"
// http://docs.neo4j.org/chunked/stable/rest-api-service-root.html
func NewConnection(uri string) (c *Connection, err error) {

	logger.SetLevel(logger.LevelDebug)

	c = &Connection{httpClient: &http.Client{}, Uri: uri}

	// prepare request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// perform request
	data, err := c.performRequest(req) // gets []byte

	// unmarshal to Connection obj
	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}
	return
}

// ExecuteCypher will return a slice of "rows", each "row" is a []*json.RawMessage representing
// a slice of node properties that the user can unmarshal themselves
func (c *Connection) ExecuteCypher(cypher string, params *map[string]interface{}) (rows [][]*json.RawMessage, err error) {

	statement := &TransactionStatement{
		Cypher:      cypher,
		Params:      *params,
		ResultTypes: []string{"ROW"},
	}

	// create a new transaction for one single statement
	// http://neo4j.com/docs/stable/rest-api-transactional.html#rest-api-begin-and-commit-a-transaction-in-one-request
	transaction := struct {
		Statements []*TransactionStatement `json:"statements"`
	}{}
	transaction.Statements = []*TransactionStatement{statement}

	// prepare request
	reqData, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	reqBuf := bytes.NewBuffer(reqData)
	uri := joinPath([]string{c.TransactionURI, "commit"})
	req, err := http.NewRequest("POST", uri, reqBuf)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json; charset=UTF-8")
	req.Header.Add("Content-Type", "application/json")

	logger.Debugf("%v: %v", cypher, *params)

	// make request
	data, err := c.performRequest(req)
	if err != nil {
		return
	}
	resp := &TransactionResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}
	if len(resp.Errors) > 0 {
		err = fmt.Errorf("%v: %v", resp.Errors[0].Code, resp.Errors[0].Message)
		return
	}
	if len(resp.Results) == 0 {
		return
	}

	// expecting only one result, since it's a single statement transaction
	rows = resp.Results[0].RowData()

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
