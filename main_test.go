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
