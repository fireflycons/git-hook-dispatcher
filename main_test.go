package main

import (
	"os"
	"path/filepath"
	"testing"
)

func mustGetHookBinary() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, "0-githooks", "hook.exe")
}

func TestProcessPreCommitHook(t *testing.T) {
	// Hook is found in hook dir for this repo
	if result, found := processHook(mustGetHookBinary(), "pre-commit"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}
}

func TestProcessPreApplyPatch(t *testing.T) {
	// Hook is found in 0-0githooks-shared
	if result, found := processHook(mustGetHookBinary(), "pre-applypatch"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}
}

func TestProcessReferenceTransaction(t *testing.T) {
	// PowerShell script can process data passed on stdin
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
	if result, found := processHook(mustGetHookBinary(), "reference-transaction"); !found {
		t.Error("Hook was not found")
	} else {
		t.Logf("returned %d", result)
	}

	if err := tmpfile.Close(); err != nil {
		t.Error(err)
	}
}

func TestMissingHookScriptExitsWithZero(t *testing.T) {
	// If a script isn't found, that should not stop git from continuing
	result, found := processHook(mustGetHookBinary(), "pre-push")

	if found {
		t.Error("Hook script pre-push should no be present in 0-gihooks.")
	}

	if result != 0 {
		t.Error("Missing hook should exit with zero")
	}
}
