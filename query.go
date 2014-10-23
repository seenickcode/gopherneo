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

	if len(resp.Errors) > 0 {
		err = errors.New(resp.Errors[0].Code + ": " + resp.Errors[0].Message)
		return
	}

	// convert each result row slice to a proper
	// slice of raw JSON the user can unmarshal later
	err = result.rawDataToRows()
	if err != nil {
		return
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
	ColumnNames []string             `json:"columns"`
	RawData     [][]*json.RawMessage `json:"data"`
	Rows        [][]byte
}

// // represents a chunk of data but in the transactional context,
// // holds the "rest" version of a result row
// type DataRow struct {
// 	Data []*json.RawMessage `json:"rest"`
// }

func (qr *QueryResult) rawDataToRows() (err error) {

	qr.Rows = make([][]byte, len(qr.RawData))

	for rowNdx, rawRow := range qr.RawData {

		m := make(map[string]*json.RawMessage)

		// pop each slice item of raw json into our map
		for colNdx, rawRowCol := range rawRow {
			colName := qr.ColumnNames[colNdx]
			m[colName] = rawRowCol
		}

		data, err := json.Marshal(m)
		if err != nil {
			return
		}

		qr.Rows[rowNdx] = json.Unmarshal(data)
	}

	// rs := make([]map[string]*json.RawMessage, len(cq.cr.Data))
	// for rowNum, row := range cq.cr.Data {
	// 	m := map[string]*json.RawMessage{}
	// 	for colNum, col := range row {
	// 		name := cq.cr.Columns[colNum]
	// 		m[name] = col
	// 	}
	// 	rs[rowNum] = m
	// }
	// b, err := json.Marshal(rs)
	// if err != nil {
	// 	logPretty(err)
	// 	return err
	// }
	// return json.Unmarshal(b, v)
}
