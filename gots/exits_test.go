package gots

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeepJSONIsRefusedNotFatal decodes and encodes JSON nested 5000 deep
// on Deno, where Go would follow it: gots lowers json's depth limit so the
// JavaScript stack never runs out on input from outside, and the program
// gets an error instead of dying.
func TestDeepJSONIsRefusedNotFatal(t *testing.T) {
	run := jsRuntime(t)
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/exits/jsondeep"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	mod := filepath.Join(out, "gots", "testdata", "exits", "jsondeep", "mod.ts")
	stdout, stderr, code := runProgram(exec.Command(run[0], append(run[1:], mod)...))
	want := "arrays refused: true\nraw refused: true\nobjects refused: true\nvalid: false\nmarshal refused: true\n"
	if code != 0 || stdout != want {
		t.Errorf("status %d, output\n%s\nwant\n%s\n%s", code, stdout, want, stderr)
	}
}

// TestProgramsEndAsGoDoes runs programs in testdata/exits, which end
// badly, under go run and Deno: the same exit status, the same output,
// and the line of stderr that says why.
func TestProgramsEndAsGoDoes(t *testing.T) {
	run := jsRuntime(t)
	for _, tc := range []struct{ dir, why string }{
		{"deadlock", "fatal error: all goroutines are asleep - deadlock!"},
		{"gopanic", "panic: boom 42"},
	} {
		t.Run(tc.dir, func(t *testing.T) {
			dir := "./testdata/exits/" + tc.dir
			goOut, goErr, goCode := runProgram(exec.Command("go", "run", dir))
			out := t.TempDir()
			if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{dir}}); err != nil {
				t.Fatalf("transpile: %v", err)
			}
			mod := filepath.Join(out, "gots", "testdata", "exits", tc.dir, "mod.ts")
			tsOut, tsErr, tsCode := runProgram(exec.Command(run[0], append(run[1:], mod)...))
			// go run reports a program's failure as its own exit status 1.
			if goCode == 1 && strings.Contains(goErr, "exit status 2") {
				goCode = 2
			}
			if tsCode != goCode || tsOut != goOut {
				t.Errorf("go: status %d, output %q\nts: status %d, output %q\n%s", goCode, goOut, tsCode, tsOut, tsErr)
			}
			for _, s := range []string{goErr, tsErr} {
				if !strings.Contains(s, tc.why) {
					t.Errorf("stderr does not say %q:\n%s", tc.why, s)
				}
			}
		})
	}
}

func runProgram(cmd *exec.Cmd) (stdout, stderr string, code int) {
	var o, e bytes.Buffer
	cmd.Stdout, cmd.Stderr = &o, &e
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		code = exit.ExitCode()
	}
	return o.String(), e.String(), code
}
