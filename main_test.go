package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessPreCommitHook(t *testing.T) {
	wd, _ := os.Getwd()

	actualBin := filepath.Join(wd, "0-githooks", "hook.exe")

	if result, found := processHook(actualBin, "pre-commit"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}
}

func TestProcessPreApplyPatch(t *testing.T) {
	wd, _ := os.Getwd()
	actualBin := filepath.Join(wd, "0-githooks", "hook.exe")

	if result, found := processHook(actualBin, "pre-applypatch"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}
}

func TestProcessReferenceTransaction(t *testing.T) {
	wd, _ := os.Getwd()
	actualBin := filepath.Join(wd, "0-githooks", "hook.exe")

	content := []byte("6c6c4afa4352441ea7b1834eac7bc70aee8248ea 37df81edea7f798982b66f4eadac531d3e730c88 HEAD\n6c6c4afa4352441ea7b1834eac7bc70aee8248ea 37df81edea7f798982b66f4eadac531d3e730c88 refs/heads/master\n")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Error(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Error(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	if result, found := processHook(actualBin, "reference-transaction"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}

	if err := tmpfile.Close(); err != nil {
		t.Error(err)
	}

}
