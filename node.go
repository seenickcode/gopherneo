package gopherneo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tideland/goas/v2/logger"
)

func (c *Connection) FindNode(label string, key string, val interface{}, result interface{}) (found bool, err error) {

	logger.Debugf("fetching %v node where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to find nodes")
		return
	}

	rows, err := c.FindNodesWithValuePaginated(label, key, val, "", 0, 0)
	if err != nil || len(rows) == 0 {
		return
	}
	if len(rows) > 1 {
		err = fmt.Errorf("found more than one %v node where '%v'='%v'", label, key, val)
		return
	}
	if len(rows) == 1 {
		row := rows[0] // []*json.RawMessage
		err = json.Unmarshal(*row[0], &result)
		found = true
	}

	return
}

func (c *Connection) FindNodesWithValuePaginated(label string, key string, val interface{}, orderClause string, pg int, pgSize int) (rows [][]*json.RawMessage, err error) {
	logger.Debugf("fetching %v nodes where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to find nodes")
		return
	}

	// determine where part
	params := &map[string]interface{}{}
	whereCypher, whereParams := cypherForWhere("n", key, val, true)
	if len(whereParams) > 0 {
		*params = whereParams
	}
	pagPart := cypherForPagination(pg, pgSize)

	// TODO ghetto. pass in order object instead of strings
	orderPart := orderClause

	parts := []string{
		fmt.Sprintf("MATCH (n:%v)", label),
		whereCypher,
		"RETURN n",
		orderPart,
		pagPart,
	}

	cypher := joinUsing(parts, " ")

	rows, err = c.ExecuteCypher(cypher, params)
	if err != nil {
		return
	}

	return
}

func (c *Connection) CreateNode(label string, props *map[string]interface{}, result interface{}) (err error) {

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

	rows, err := c.ExecuteCypher(cypher, params)
	if err != nil {
		return
	}
	if result != nil && len(rows) == 1 {
		row := rows[0] // []*json.RawMessage
		err = json.Unmarshal(*row[0], &result)
	}

	return
}

func (c *Connection) UpdateNode(label string, key string, val interface{}, props *map[string]interface{}, result interface{}) (err error) {

	logger.Debugf("updating %v node with: %v", label, *props)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to update a node")
		return
	}

	params := make(map[string]interface{})

	// normally we'd use 'SET {props}' but that replaces _all_ the node's props
	// and here we want to do it with only the props the user provides
	setCypher, setParams := cypherForSetProps("n", props)
	// copy params
	for k, v := range setParams {
		params[k] = v
	}

	whereCypher, whereParams := cypherForWhere("n", key, val, true)
	// copy params
	for k, v := range whereParams {
		params[k] = v
	}

	parts := []string{
		fmt.Sprintf("MATCH (n:%v)", label),
		whereCypher,
		setCypher,
		"RETURN n",
	}

	cypher := joinUsing(parts, " ")

	rows, err := c.ExecuteCypher(cypher, &params)
	if err != nil {
		return
	}
	if result != nil && len(rows) == 1 {
		row := rows[0] // []*json.RawMessage
		err = json.Unmarshal(*row[0], &result)
	}

	return
}

func (c *Connection) DeleteNodes(label string, key string, val interface{}) (err error) {

	logger.Debugf("deleting %v node where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to delete nodes")
		return
	}

	params := &map[string]interface{}{}
	whereCypher, whereParams := cypherForWhere("n", key, val, true)
	if len(whereParams) > 0 {
		*params = whereParams
	}

	parts := []string{
		fmt.Sprintf("MATCH (n:%v)", label),
		whereCypher,
		"OPTIONAL MATCH (n)-[r]-() DELETE n, r",
	}

	cypher := joinUsing(parts, " ")

	_, err = c.ExecuteCypher(cypher, params)

	return
}

// LinkNodes requires that keys and values provided have unique properties for those fields
func (c *Connection) LinkNodes(label1 string, key1 string, val1 string, label2 string, key2 string, val2 string, relName string, relProps *map[string]interface{}, relStruct interface{}) (err error) {

	logger.Debugf("linking %v node where '%v'='%v' to %v node where '%v'='%v'", label1, key1, val1, label2, key2, val2)

	if len(label1) == 0 || len(label2) == 0 {
		err = fmt.Errorf("labels are required to link nodes")
		return
	}
	if len(relName) == 0 {
		err = fmt.Errorf("relName is required to link nodes")
		return
	}

	cypher := fmt.Sprintf(`
      MATCH (n1:%v), (n2:%v)
      WHERE n1.%v={val1} AND n2.%v={val2}
      CREATE UNIQUE (n1)-[r:%v {relProps}]->(n2)
      RETURN r`, label1, label2, key1, key2, relName)

	params := &map[string]interface{}{
		"val1":     val1,
		"val2":     val2,
		"relProps": relProps,
	}

	rows, err := c.ExecuteCypher(cypher, params)

	if err == nil {
		row := rows[0] // []*json.RawMessage
		err = json.Unmarshal(*row[0], &relStruct)
	}

	return
}

func cypherForSetProps(alias string, props *map[string]interface{}) (cypher string, params map[string]interface{}) {

	params = make(map[string]interface{})

	parts := make([]string, len(*props))
	i := 0
	for key, val := range *props {
		paramKey := fmt.Sprintf("setval%d", i)
		parts[i] = fmt.Sprintf("n.%v={%v}", key, paramKey)
		params[paramKey] = val
		i++
	}
	cypher = "SET " + joinUsing(parts, ", ")
	return
}

func cypherForWhere(alias string, key string, val interface{}, inclKeyword bool) (cypher string, params map[string]interface{}) {
	var b bytes.Buffer
	if len(key) > 0 {
		if inclKeyword {
			b.WriteString("WHERE ")
		}
		b.WriteString(fmt.Sprintf("%v.%v={whereval}", alias, key))
		params = make(map[string]interface{})
		params["whereval"] = val
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
