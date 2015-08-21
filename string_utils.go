package main

import (
	"strings"
)

func breakLongLine(input string, lineLen int) string {
	fields := strings.Fields(input)
	output := ""
	tmpLen := 0

	for _, field := range fields {
		fieldLen := len(field)
		tmpLen += fieldLen + 1
		if tmpLen >= lineLen {
			output += "\n"
			tmpLen = fieldLen
		}
		output += field + " "
	}

	return output
}

func indentMultilines(input string, indentLen int) string {
	return strings.Repeat(" ", indentLen) + strings.Replace(input, "\n", "\n  ", -1)
}
