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
func (c *Connection) Query(statement *Query) (rows []map[string]interface{}, err error) {

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

	// unmarshal our QueryResponse
	resp := &QueryResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}
	if len(resp.Errors) > 0 {
		err = errors.New(resp.Errors[0].Code + ": " + resp.Errors[0].Message)
		return
	}

	// since we're only expecting one result (single transaction),
	// convert it to a []map[string]interface{} using the interface we passed in
	if len(resp.Results) == 0 {
		return
	}

	result := resp.Results[0]
	rows = make([]map[string]interface{}, len(result.Rows))

	for rowNdx, rawRow := range result.Rows {
		m := make(map[string]interface{})
		for colNdx, colValue := range rawRow.Cols {
			colName := result.ColumnNames[colNdx]
			m[colName] = colValue
		}
		rows[rowNdx] = m
	}
	return
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
	Cols []interface{} `json:"rest"`
}
