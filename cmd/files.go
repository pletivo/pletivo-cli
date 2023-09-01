package cmd

import (
	"fmt"
	"github.com/pletivo/pletivo-cli/internal/pletivo"
	"github.com/spf13/cobra"
	"path/filepath"
)

func FilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "files [DIR]",
		Short: "Renders files graph",
		Long: `The files command iterates the directory structure
and renders Cypher Script to make the tree of files and directories 
inside neo4j. For example:

pletivo-cli files ../work/some-project
`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := filepath.WalkDir(args[0], pletivo.DirVisitor(args[0]))
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
}

// filesCmd represents the files command
var filesCmd = FilesCmd()

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
