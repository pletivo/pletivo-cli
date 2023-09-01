package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandReturnsCypherScripts(t *testing.T) {
	cmd := FilesCmd()
	pwd, _ := os.Getwd()
	pwd = filepath.Join(pwd, "..")
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{pwd})
	err := cmd.Execute()

	if err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), "MATCH") {
		t.Failed()
	}
}
