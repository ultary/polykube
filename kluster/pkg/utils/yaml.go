package utils

import (
	"bytes"
)

func SplitYAML(data []byte) [][]byte {
	var parts [][]byte
	docs := bytes.Split(data, []byte("\n---\n"))
	for _, doc := range docs {
		trimmed := bytes.TrimSpace(doc)
		if len(trimmed) > 0 {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
