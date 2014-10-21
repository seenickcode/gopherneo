package gopherneo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CypherResponse struct {
	Columns []string             `json:"columns"`
	Data    [][]*json.RawMessage `json:"data"`
}

// perform a cypher query
// http://docs.neo4j.org/chunked/stable/rest-api-cypher.html#rest-api-send-a-query
func (c *Connection) Query(query string) (cypResp CypherResponse, err error) {
	return c.QueryWithParams(query, nil)
}

// perform a cypher query with params
// http://docs.neo4j.org/chunked/stable/rest-api-cypher.html#rest-api-use-parameters
func (c *Connection) QueryWithParams(query string, params map[string]string) (cypResp CypherResponse, err error) {

	log.Print(query)

	// compose URI
	uri := c.CypherURI

	// prepare request body
	cypherReq := cypherRequest{Query: query, Params: params}
	reqData, err := json.Marshal(cypherReq)
	if err != nil {
		return
	}
	reqBuf := bytes.NewBuffer(reqData)

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

	// deserialize
	err = json.Unmarshal(data, &cypResp)
	if err != nil {
		return
	}
	return

}

// private types

type cypherRequest struct {
	Query  string            `json:"query"`
	Params map[string]string `json:"params"`
}
