package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Parent interface {
	Dir() string
}

type ParentDir struct {
	dir string
}

func (p *ParentDir) Dir() string {
	return p.dir
}

func parentDir(path string) Parent {
	base := filepath.Dir(path)
	if base == path {
		return nil
	}

	return &ParentDir{dir: base}
}

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}

	if strings.HasPrefix(path, ".") {
		return nil
	}

	if strings.HasPrefix(path, "tmp") {
		return nil
	}

	if strings.Contains(path, ".DS_Store") {
		return nil
	}

	if info.IsDir() {
		fmt.Printf("CREATE (:SourceDir:SourceLayer{path: %s});\n", strconv.Quote(path))
		parent := parentDir(path)

		if parent != nil {
			fmt.Printf("MATCH (pd:SourceDir), (cd:SourceDir) WHERE pd.path = %s AND cd.path = %s CREATE (cd)-[:SEATS_IN]->(pd);\n", strconv.Quote(parent.Dir()), strconv.Quote(path))
		}
	} else {
		fmt.Printf("CREATE (:SourceFile:SourceLayer{filename: %s});\n", strconv.Quote(path))
		parent := parentDir(path)

		if parent != nil {
			fmt.Printf("MATCH (pd:SourceDir), (f:SourceFile) WHERE pd.path = %s AND f.filename = %s CREATE (f)-[:SEATS_IN]->(pd);\n", strconv.Quote(parent.Dir()), strconv.Quote(path))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <directory>")
		return
	}

	root := os.Args[1]
	err := filepath.Walk(root, visit)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
