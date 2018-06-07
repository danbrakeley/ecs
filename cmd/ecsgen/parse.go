package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func getComponents(r io.Reader) []Component {
	var coms []Component

	reComBase := regexp.MustCompile(`^\s*ecs\.ComponentBase(\s+.*|$)`)
	reTypeStruct := regexp.MustCompile(`^\s*type\s+(\w+)\s+struct.*$`)
	reSerializer := regexp.MustCompile(`^\s*func\s+\(\s*\w+\s+\*(\w+)\s*\)\s+Serialize\(.*$`)
	reDeserializer := regexp.MustCompile(`^\s*func\s+Deserialize(\w+)\(.*$`)
	line := 0
	lastLine := ""
	hasSerializer := make(map[string]bool)
	hasDeserializer := make(map[string]bool)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		curLine := scanner.Text()
		if reComBase.MatchString(curLine) {
			m := reTypeStruct.FindStringSubmatch(lastLine)
			if len(m) > 1 && len(m[1]) > 0 {
				coms = append(coms, Component{Line: line, Name: m[1]})
			}
		}

		// if this is a Serializer, note the Component name for later
		m := reSerializer.FindStringSubmatch(curLine)
		if len(m) > 1 && len(m[1]) > 0 {
			hasSerializer[m[1]] = true
		}

		// if this is a Deserializer, note the Component name for later
		m = reDeserializer.FindStringSubmatch(curLine)
		if len(m) > 1 && len(m[1]) > 0 {
			hasDeserializer[m[1]] = true
		}

		lastLine = curLine
		line++
	}

	// fill in deserializer info for any found components
	for i := range coms {
		coms[i].HasSerializer = hasSerializer[coms[i].Name]
		coms[i].HasDeserializer = hasDeserializer[coms[i].Name]
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return coms
}
