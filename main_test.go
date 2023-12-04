package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func mustGetHookBinary() string {
	// "mock" location of hook binary
	// such that githooks directory is discovered.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, "0-githooks", "hook.exe")
}

func vaildateScriptResult(t *testing.T, hookName string, expectedResult int, shouldBeFound bool) {
	result, found := processHook(mustGetHookBinary(), hookName)

	require.Equal(t, shouldBeFound, found, func() string {
		if shouldBeFound && !found {
			return fmt.Sprintf("%s.ps1 was not found when it should have been.", hookName)
		}
		return fmt.Sprintf("%s.ps1 was found when it should not have been.", hookName)
	}())

	require.Equalf(t, expectedResult, result, "Script exited with %d when %d was expected.", result, expectedResult)
}

func TestProcessPreCommitHookFoundInRepoHooks(t *testing.T) {
	// Hook is found in hook dir for this repo.
	vaildateScriptResult(t, "pre-commit", 0, true)
}

func TestProcessPreApplyPatchFoundInSharedHooks(t *testing.T) {
	// Hook is found in 0-githooks-shared.
	vaildateScriptResult(t, "pre-applypatch", 0, true)
}

func TestProcessReferenceTransactionHandlesStdin(t *testing.T) {
	// PowerShell script can process data passed on stdin.
	//
	// Wire up a temp file to act as process stdin.
	tmpfile, err := os.CreateTemp("", "mock-stdin")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	content := []byte("6c6c4afa4352441ea7b1834eac7bc70aee8248ea 37df81edea7f798982b66f4eadac531d3e730c88 HEAD\n6c6c4afa4352441ea7b1834eac7bc70aee8248ea 37df81edea7f798982b66f4eadac531d3e730c88 refs/heads/master\n")
	if _, err := tmpfile.Write(content); err != nil {
		t.Error(err)
	}

	if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
		t.Error(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin.

	os.Stdin = tmpfile

	vaildateScriptResult(t, "reference-transaction", 0, true)

	if err := tmpfile.Close(); err != nil {
		t.Error(err)
	}
}

func TestMissingHookScriptExitsWithZero(t *testing.T) {
	// If a script isn't found, that should not stop git from continuing
	vaildateScriptResult(t, "pre-push", 0, false)
}

func TestPowerShellThrowReturnsNonZero(t *testing.T) {
	// If a script (in this case pre-rebase) throws an exception from PowerShell,
	// then it should stop git by returning 1.
	vaildateScriptResult(t, "pre-rebase", 1, true)
}
