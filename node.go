package gopherneo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tideland/goas/v2/logger"
)

func (c *Connection) CreateNodeWithLabel(label string, props *map[string]interface{}, result interface{}) (err error) {

	logger.Debugf("creating %v node with: %v", label, *props)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to create a node")
		return
	}

	cypher := fmt.Sprintf(`CREATE (n:%v {p}) RETURN n`, label)

	// add my cypher props to a map[string]interface{}
	params := &map[string]interface{}{
		"p": props,
	}

	rows, err := c.Query(cypher, params)
	if err != nil {
		return
	}
	if len(rows) != 1 {
		err = fmt.Errorf("couldn't create node with %v, expected only 1 node", props)
		return
	}
	row := rows[0] // []*json.RawMessage

	// convert our single transaction result
	if result != nil {
		err = json.Unmarshal(*row[0], &result)
	}

	return
}

func (c *Connection) FindNodesWithValuePaginated(label string, key string, val string, pg int, pgSize int) (rows [][]*json.RawMessage, err error) {

	logger.Debugf("fetching %v nodes where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to find nodes")
		return
	}

	// determine where part
	params := &map[string]interface{}{}
	wherePart, whereParams := cypherForWhere("n", key, val, true)
	if len(*whereParams) > 0 {
		params = whereParams
	}
	pagPart := cypherForPagination(pg, pgSize)

	parts := []string{
		fmt.Sprintf("MATCH (n:%v)", label),
		wherePart,
		"RETURN n",
		pagPart,
	}

	cypher := joinUsing(parts, " ")

	rows, err = c.Query(cypher, params)
	if err != nil {
		return
	}

	return
}

func (c *Connection) DeleteNodes(label string, key string, val string) (err error) {

	logger.Debugf("deleting %v node where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to delete nodes")
		return
	}

	// determine where part
	params := &map[string]interface{}{}
	wherePart, whereParams := cypherForWhere("n", key, val, true)
	if len(*whereParams) > 0 {
		params = whereParams
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("MATCH (n:%v) ", label))
	buffer.WriteString(wherePart)
	buffer.WriteString("OPTIONAL MATCH (n)-[r]-() DELETE n, r")

	cypher := buffer.String()

	_, err = c.Query(cypher, params)

	return
}

func cypherForWhere(alias string, key string, val string, inclKeyword bool) (cypher string, params *map[string]interface{}) {
	params = &map[string]interface{}{}

	var b bytes.Buffer
	if len(key) > 0 && len(val) > 0 {
		if inclKeyword {
			b.WriteString("WHERE ")
		}
		b.WriteString(fmt.Sprintf("%v.%v={val}", alias, key))
		params = &map[string]interface{}{
			"val": val,
		}
	}
	cypher = b.String()
	return
}

func cypherForPagination(pg int, pgSize int) (cypher string) {

	if pg < 0 {
		pg = 0
	}
	if pgSize > 0 {
		skip := pg * pgSize
		cypher = cypher + fmt.Sprintf("SKIP %d LIMIT %d", skip, pgSize)
	}
	return
}
