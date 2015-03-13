package gopherneo

import "regexp"

func EscapeStringForCypherRegex(in string) string {
	//fmt.Printf("in: %v\n", in)
	r1 := regexp.MustCompile("(\\')")
	r2 := regexp.MustCompile("(\\(|\\))")
	r3 := regexp.MustCompile("(&)")
	r4 := regexp.MustCompile("(\\*)")
	r5 := regexp.MustCompile("(\")")
	r6 := regexp.MustCompile("(\\+)")
	r7 := regexp.MustCompile("(\\$)") // TODO test and add all remaining regex meta characters!

	out := in
	out = r1.ReplaceAllString(out, "\\\\$1")
	out = r2.ReplaceAllString(out, "\\\\$1")
	out = r3.ReplaceAllString(out, "\\\\$1")
	out = r4.ReplaceAllString(out, "\\\\$1")
	out = r5.ReplaceAllString(out, "\\$1")
	out = r6.ReplaceAllString(out, "\\\\$1")
	out = r7.ReplaceAllString(out, "\\\\$1")

	//fmt.Printf("out: %v\n", out)
	return out
}
