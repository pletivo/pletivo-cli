/*
Copyright © 2023 Anatolii Lapshin <holywarez@gmail.com>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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

func visit(rootDir string) fs.WalkDirFunc {
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

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files [DIR]",
	Short: "Renders files graph",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := filepath.WalkDir(args[0], visit(args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
