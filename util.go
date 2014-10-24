package gopherneo

import (
	"bytes"
	"strconv"
	"testing"
	"time"
	"runtime"
)

// string utils
//

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

// generators
//

func generateTimestamp() int64 {
	return time.Now().UnixNano()
}

// testing
//

func assertOk(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		lineNo := strconv.Itoa(line)
		t.Errorf(file + ":" + lineNo + ": %# v\n", err)
		t.Error(err)
	}
}

// experimental
//

func cypherizeNode(label string, props map[string]interface{}, alias string) (result string) {

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
		// 	value, _ = strconv.ParseFloat(v.(string), 64)
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
