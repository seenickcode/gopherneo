package gopherneo

import (
	"encoding/json"
	"fmt"
)

func (c *Connection) CreateNodeWithLabel(label string, props *map[string]interface{}, result interface{}) (err error) {

	cypher := fmt.Sprintf(`
    CREATE (n:%v {p}) RETURN n
  `, label)

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
	row := rows[0]

	err = json.Unmarshal(*row, &result)

	return
}
