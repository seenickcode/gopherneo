package gopherneo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Query struct {
	Cypher      string                 `json:"statement"`
	Params      map[string]interface{} `json:"parameters"`
	ResultTypes []string               `json:"resultDataContents"`
}

// NewQuery constructs a new query that the Neo4j transactional
// endpoint recognizes
func NewQuery(cypher string, params map[string]interface{}) *Query {
	qs := &Query{}
	qs.Cypher = cypher
	qs.Params = make(map[string]interface{})
	qs.ResultTypes = []string{"REST"}
	return qs
}

// Query leverages the official Neo4j transactional endpoint
// and commits a single statement immediately.
func (c *Connection) Query(statement *Query) (resp QueryResponse, err error) {

	uri := joinPath([]string{c.TransactionURI, "commit"})

	// create a new transaction
	// initialized with one single statement
	transaction := struct {
		Statements []*Query `json:"statements"`
	}{}
	transaction.Statements = []*Query{statement}

	// prepare request
	reqData, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	reqBuf := bytes.NewBuffer(reqData)
	req, err := http.NewRequest("POST", uri, reqBuf)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json; charset=UTF-8")
	req.Header.Add("Content-Type", "application/json")

	// make request
	data, err := c.performRequest(req)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}

	return
}

// represents a transactional response
type QueryResponse struct {
	Results []QueryResult `json:"results"`
	Errors  []interface{} `json:"errors"`
}

// represents a transactional response result
type QueryResult struct {
	ColumnNames []string  `json:"columns"`
	Rows        []DataRow `json:"data"`
}

// represents a chunk of data but in the transactional context,
// holds the "rest" version of a result row
type DataRow struct {
	Data []interface{} `json:"rest"`
}
