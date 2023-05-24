package main

import (
	"encoding/json"
	"os"
	alignment "tree_alignment/internal"
)

func readContent(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

func LoadGraphDescription(graphPath string) []alignment.NodeDescription {
	content := readContent(graphPath)
	var description []alignment.NodeDescription
	if err := json.Unmarshal(content, &description); err != nil {
		panic(err)
	}
	return description
}
