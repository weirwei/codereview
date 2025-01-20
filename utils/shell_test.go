package utils

import "testing"

func TestShell(t *testing.T) {
	t.Run("git diff", func(t *testing.T) {
		result, err := ShellExec("git", "diff", "origin/main...HEAD", "--name-only", "--diff-filter=d")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	})
}
