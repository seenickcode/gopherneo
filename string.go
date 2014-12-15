package gopherneo

import (
	"bytes"
	"strconv"
)

func joinPath(strings []string) string {
	return joinUsing(strings, "/")
}

func join(strings []string) string {
	var buffer bytes.Buffer
	for _, s := range strings {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func joinUsing(strings []string, delimiter string) string {
	var buffer bytes.Buffer
	for ndx, s := range strings {
		buffer.WriteString(s)
		if ndx != len(strings)-1 {
			buffer.WriteString(delimiter)
		}
	}
	return buffer.String()
}

func propsToCypherString(label string, props map[string]interface{}, alias string) (result string) {

	// create a slice of keys, values
	propParts := []string{}
	for k, v := range props {
		key := k
		value := ""
		switch v.(type) {
		case int:
			value = strconv.Itoa(v.(int))
		// TODO support floats. ** figure out how to convert float to string **
		// case float64:
		//  value, _ = strconv.ParseFloat(v.(string), 64)
		case string:
			value = ("'" + v.(string) + "'")
		default:
			value = ""
		}
		if len(value) > 0 {
			propParts = append(propParts, key+": "+value)
		}
	}

	// construct our cypher
	result = "(" + alias + ":" + label + " { " + joinUsing(propParts, ", ") + " })"

	return
}
