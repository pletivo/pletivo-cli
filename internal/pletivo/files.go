package pletivo

import (
	"fmt"
	"io/fs"
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

func DirVisitor(rootDir string) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		path, _ = filepath.Rel(rootDir, path)

		if strings.HasPrefix(path, ".") {
			return nil
		}

		if strings.HasPrefix(path, "tmp") {
			return nil
		}

		if strings.Contains(path, ".DS_Store") {
			return nil
		}

		if d.IsDir() {
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
}
