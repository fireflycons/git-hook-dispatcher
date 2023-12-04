package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type scriptType int

const (
	none scriptType = iota
	powershell
	cmdOrBat
)

var (
	kernel                = syscall.MustLoadDLL("kernel32.dll")
	getModuleFileNameProc = kernel.MustFindProc("GetModuleFileNameW")
)

// Gets full path to this actual program
// from Windows kernel.
func getModuleFileName() (string, error) {
	var n uint32

	// Buffer to receive exe's path
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	r0, _, e1 := getModuleFileNameProc.Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(size))
	n = uint32(r0)

	if n == 0 {
		return "", e1
	}

	return string(utf16.Decode(b[0:n])), nil
}

func debugEnabled() bool {
	return os.Getenv("GITHOOK_DEBUG") != ""
}

func scriptNameWithoutExtension(hookDir, hookName, repo string) string {
	return filepath.Join(hookDir, repo, hookName)
}

func getHookScript(hookDir, hookName, repo string) (string, scriptType) {

	scriptMap := map[string]scriptType{
		"bat": cmdOrBat,
		"cmd": cmdOrBat,
		"ps1": powershell,
	}

	for _, dir := range []string{repo, "00-githooks-shared"} {
		for extension, typ := range scriptMap {
			script := fmt.Sprintf("%s.%s", scriptNameWithoutExtension(hookDir, hookName, dir), extension)
			if _, err := os.Stat(script); err == nil {
				if debugEnabled() {
					fmt.Printf("Will execute %s\n", script)
				}
				return script, typ
			}
		}
	}

	return "", none
}

const defaultFailedCode = 1

func getExitCode(ps *os.ProcessState) int {
	if ws, ok := ps.Sys().(syscall.WaitStatus); !ok {
		return defaultFailedCode
	} else {
		return ws.ExitStatus()
	}
}

func RunCommand(cmd *exec.Cmd) (exitCode int) {
	// Redirect pipes of script to run
	// attaching them to this process's pipes
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = getExitCode(exitError.ProcessState)
		} else {
			exitCode = defaultFailedCode
			fmt.Println(err.Error())
		}
	} else {
		// success, exitCode should be 0 if go is ok
		exitCode = getExitCode(cmd.ProcessState)
	}

	return
}

func processHook(actualBin string, hookName string) (int, bool) {

	// Hooks directory is the deirectory where this executable is found.
	hookDir := filepath.Dir(actualBin)
	hookArgs := os.Args[1:]
	wd, _ := os.Getwd()

	// Git executes the hook with the working directory set to the root of the repo.
	repoName := filepath.Base(wd)

	var cmd *exec.Cmd
	script, typ := getHookScript(hookDir, hookName, repoName)

	switch typ {
	case cmdOrBat:
		cmd = exec.Command("cmd", append([]string{"/c", script}, hookArgs...)...)
	case powershell:
		cmd = exec.Command("powershell", append([]string{"-File", script}, hookArgs...)...)
	default:
		if debugEnabled() {
			fmt.Printf("(No script found for hook %s)\n", hookName)
		}

		// Exit zero (don't stop action), and false to indicate no script run (for tests).
		return 0, false
	}

	// Exit code of script, and true to indicate script was run.
	return RunCommand(cmd), true
}

func main() {

	// First argument is the path used to invoke the program.
	// If invoked via a symlink, then the symlink name will
	// be the basename of this path, and we use this to
	// identify the actual hook.
	hookName := filepath.Base(os.Args[0])

	// This will be the absolute path to the real executable,
	// not any relative or symlinked path.
	actualBin, err := getModuleFileName()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	i := strings.LastIndex(hookName, ".")
	if i > 1 {
		// Strip extension.
		hookName = hookName[:i]
	}

	exitCode, _ := processHook(actualBin, hookName)
	os.Exit(exitCode)
}
