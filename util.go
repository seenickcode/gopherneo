package gopherneo

import (
	"bytes"
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
