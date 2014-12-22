package gopherneo

import (
	"encoding/json"
	"fmt"

	"github.com/tideland/goas/v2/logger"
)

func (c *Connection) LinkNodes(label1 string, key1 string, val1 string, label2 string, key2 string, val2 string, relName string, relProps *map[string]interface{}, resultRel interface{}) (err error) {

	logger.Debugf("linking %v node where '%v'='%v' to %v node where '%v'='%v'", label1, key1, val1, label2, key2, val2)

	if len(label1) == 0 || len(label2) == 0 {
		err = fmt.Errorf("labels are required to link nodes")
		return
	}
	if len(relName) == 0 {
		err = fmt.Errorf("relName is required to link nodes")
		return
	}
	if relProps == nil {
		relProps = &map[string]interface{}{}
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
	cr, err := c.ExecuteCypher(cypher, params)
	if err != nil {
		return
	}
	if resultRel != nil && len(cr.Rows) > 0 {
		err = json.Unmarshal(cr.Rows[0], &resultRel)
	}

	return
}

// UnlinkAllNodes removes all relationships of a specific type from a specified node
func (c *Connection) UnlinkAllNodes(label1 string, key1 string, val1 string, relName string, label2 string) (err error) {

	logger.Debugf("deleting all %v rels to %v nodes from %v node where '%v'='%v' ", relName, label2, label1, key1, val1)

	if len(label1) == 0 || len(label2) == 0 {
		err = fmt.Errorf("labels are required to unlink nodes")
		return
	}
	if len(relName) == 0 {
		err = fmt.Errorf("relName is required to unlink nodes")
		return
	}

	cypher := fmt.Sprintf(`
      MATCH (n1:%v)-[r:%v]-(n2:%v)
      WHERE n1.%v={val1} 
      DELETE r`, label1, relName, label2, key1)

	params := &map[string]interface{}{
		"val1": val1,
	}
	_, err = c.ExecuteCypher(cypher, params)
	if err != nil {
		return
	}
	return
}

func (c *Connection) FindAllRelNodesPaginated(label string, key string, val interface{}, relLabel string, relName string, orderClause string, pg int, pgSize int) (cr CypherResult, err error) {
	logger.Debugf("fetching %v nodes where '%v'='%v'", label, key, val)

	if len(label) == 0 {
		err = fmt.Errorf("a label is required to find nodes")
		return
	}
	if len(relName) == 0 {
		err = fmt.Errorf("relName is required to unlink nodes")
		return
	}

	params := &map[string]interface{}{
		"relName": relName,
	}
	whereCypher, whereParams := cypherForWhere("n1", key, val, true)
	if len(whereParams) > 0 {
		*params = whereParams
	}
	pagPart := cypherForPagination(pg, pgSize)

	orderPart := orderClause

	parts := []string{
		fmt.Sprintf("MATCH (n1:%v)-[:%v]->(n2:%v)", label, relName, relLabel),
		whereCypher,
		"RETURN n2",
		orderPart,
		pagPart,
	}

	cypher := joinUsing(parts, " ")

	cr, err = c.ExecuteCypher(cypher, params)
	if err != nil {
		return
	}

	return
}
