package gopherneo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Query struct {
	Cypher      string                 `json:"statement"`
	Params      map[string]interface{} `json:"parameters"`
	ResultTypes []string               `json:"resultDataContents"`
}

// represents a transactional response
type QueryResponse struct {
	Results []QueryResult `json:"results"`
	Errors  []QueryError  `json:"errors"`
}

type QueryError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// represents a transactional response result
type QueryResult struct {
	ColumnNames []string    `json:"columns"`
	Rows        []ResultRow `json:"data"`
}

// holds the "rest" version of a result row
type ResultRow struct {
	Cols []*json.RawMessage `json:"rest"`
}

// NewQuery constructs a new query that the Neo4j transactional
// endpoint recognizes
func NewQuery(cypher string) *Query {
	qs := &Query{}
	qs.Cypher = cypher
	qs.Params = make(map[string]interface{})
	qs.ResultTypes = []string{"REST"}
	return qs
}

// Query leverages the official Neo4j transactional endpoint, committing
// a single statement immediately and returning a single result
func (c *Connection) Query(statement *Query) (rows []*map[string]interface{}, err error) {

	uri := joinPath([]string{c.TransactionURI, "commit"})

	// create a new transaction for one single statement
	// http://neo4j.com/docs/stable/rest-api-transactional.html#rest-api-begin-and-commit-a-transaction-in-one-request
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
	// add headers
	req.Header.Add("Accept", "application/json; charset=UTF-8")
	req.Header.Add("Content-Type", "application/json")

	// make request
	data, err := c.performRequest(req)
	if err != nil {
		return
	}

	// unmarshal our QueryResponse
	resp := &QueryResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}

	// handle error messages
	if len(resp.Errors) > 0 {
		err = errors.New(resp.Errors[0].Code + ": " + resp.Errors[0].Message)
		return
	}
	if len(resp.Results) == 0 {
		return
	}

	// deliberately passed single transaction statement, expecting only one result
	result := resp.Results[0]

	rows = make([]*map[string]interface{}, len(result.Rows))

	for ri, row := range result.Rows {
		rm := make(map[string]interface{}) // row map
		for ci, colVal := range row.Cols {
			n := result.ColumnNames[ci]
			rm[n] = colVal
		}
		rows[ri] = &rm
	}

	return
}
