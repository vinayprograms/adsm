package diagram

import "strings"

func appendLine(document []string, tabs int, line string) []string {
	return append(document, genrateTabs(tabs) + line)
}

func appendLineSpacer(document []string) []string {
	return append(document, "")
}

func genrateTabs(tabCount int) (result string) {
	for i := 0; i < tabCount; i++ {
		result += "  "
	}
	return
}

func GenerateID(s string) string {
	id := cleanup(s)
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return id
}

// Replace symbols for use in ID generator
func cleanup(str string) (cleanedStr string) {
	cleanedStr = str
	// Symbols to remove
	for _, s := range []string{".", "(", ")", "[", "]", "{", "}", "'", "`", "+", "?", ",", ":"} {
		cleanedStr = strings.ReplaceAll(cleanedStr, s, "")
	}
	// Replace with alternate
	replacements := map[string]string{
		"<":  "lt",
		">":  "gt",
		"=":  "eq",
		"\"": "\\\"",
	}
	for k, v := range replacements {
		cleanedStr = strings.ReplaceAll(cleanedStr, k, v)
	}

	return cleanedStr
}

// Wrap long lines exceeding 15 chars, along word boundaries
func wrap(s string) (wrapString string) {
	temp := strings.Fields(s)
	length := 0
	for _, str := range temp {
		length += len(str)
		if length > 15 {
			wrapString += "\\n" + str
			length = 0
		} else {
			wrapString += " " + str
		}
	}
	return strings.TrimSpace(wrapString)
}

// Wrap long lines, but use HTML line breaks.
func htmlwrap(s string) (wrapString string) {
	temp := strings.Fields(s)
	length := 0
	for _, str := range temp {
		length += len(str)
		if length > 15 {
			wrapString += "<br></br>" + str
			length = 0
		} else {
			wrapString += " " + str
		}
	}
	return strings.TrimSpace(wrapString)
}

// Check array membership
func contains[T comparable](item T, array []T) bool {
	for _, x := range array {
		if item == x {
			return true
		}
	}

	return false
}