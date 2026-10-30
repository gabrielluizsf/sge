package gots

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestDeepJSONIsRefusedNotFatal decodes and encodes JSON nested 5000 deep
// on Deno, where Go would follow it: gots lowers json's depth limit so the
// JavaScript stack never runs out on input from outside, and the program
// gets an error instead of dying.
func TestDeepJSONIsRefusedNotFatal(t *testing.T) {
	t.Parallel()
	run := jsRuntime(t)
	out := programs(t)
	mod := filepath.Join(out, "gots", "testdata", "exits", "jsondeep", "mod.ts")
	stdout, stderr, code := runProgram(exec.Command(run[0], append(run[1:], mod)...))
	want := "arrays refused: true\nraw refused: true\nobjects refused: true\nvalid: false\nmarshal refused: true\n"
	if code != 0 || stdout != want {
		t.Errorf("status %d, output\n%s\nwant\n%s\n%s", code, stdout, want, stderr)
	}
}

// TestDeepXMLIsRefusedNotFatal is TestDeepJSONIsRefusedNotFatal for
// encoding/xml, whose own depth limit is sized for a stack that grows.
func TestDeepXMLIsRefusedNotFatal(t *testing.T) {
	t.Parallel()
	run := jsRuntime(t)
	out := programs(t)
	mod := filepath.Join(out, "gots", "testdata", "exits", "xmldeep", "mod.ts")
	stdout, stderr, code := runProgram(exec.Command(run[0], append(run[1:], mod)...))
	want := "100 true true\n300 false true\n4000 false true\n20000 false true\n"
	if code != 0 || stdout != want {
		t.Errorf("status %d, output\n%s\nwant\n%s\n%s", code, stdout, want, stderr)
	}
}

// TestDeepGobIsRefusedNotFatal decodes, on Deno, gob streams nested 5000
// deep, which the decoder walks by recursion, the stack exhaustion of
// CVE-2024-34156: a type nested so, which Go limits to 10000 levels, a
// limit sized for a stack that grows, and data nested so in a recursive
// type, which Go does not limit at all. Both are refused with an error
// rather than end the program.
func TestDeepGobIsRefusedNotFatal(t *testing.T) {
	t.Parallel()
	run := jsRuntime(t)
	// A type nested 5000 deep.
	nested := reflect.TypeFor[int]()
	for range 5000 {
		nested = reflect.StructOf([]reflect.StructField{{Name: "F", Type: nested}})
	}
	var deepType bytes.Buffer
	if err := gob.NewEncoder(&deepType).Encode(reflect.New(nested).Interface()); err != nil {
		t.Fatal(err)
	}
	// A recursive type, one level of it, holding data nested 5000 deep:
	// the type passes any limit on types, the data has to be held to one.
	type S struct{ F *S }
	type W struct {
		Hello int
		X     *S
	}
	list := &S{}
	for range 5000 {
		list = &S{F: list}
	}
	var deepData bytes.Buffer
	if err := gob.NewEncoder(&deepData).Encode(W{Hello: 1, X: list}); err != nil {
		t.Fatal(err)
	}
	out := programs(t)
	mod := filepath.Join(out, "gots", "testdata", "exits", "gobdeep", "mod.ts")
	for name, b := range map[string][]byte{"type": deepType.Bytes(), "data": deepData.Bytes()} {
		in := filepath.Join(t.TempDir(), name+".gob")
		if err := os.WriteFile(in, b, 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(run[0], append(run[1:], mod)...)
		cmd.Env = append(os.Environ(), "GOB_INPUT="+in)
		stdout, stderr, code := runProgram(cmd)
		if code != 0 || stdout != "refused: true\n" {
			t.Errorf("%s nested deep: status %d, output %q\n%s", name, code, stdout, stderr)
		}
	}
}

// TestHugeContentLengthIsRefused sends a server, and hands a client, a
// Content-Length past 2^53, which a JavaScript number would round: both
// refuse it on Deno rather than read another length than was meant, as
// request smuggling needs.
func TestHugeContentLengthIsRefused(t *testing.T) {
	t.Parallel()
	run := jsRuntime(t)
	out := programs(t)
	mod := filepath.Join(out, "gots", "testdata", "exits", "hugelength", "mod.ts")
	stdout, stderr, code := runProgram(exec.Command(run[0], append(run[1:], mod)...))
	want := "server: \"HTTP/1.1 400 Bad Request\\r\\n\"\nclient refused: true\n"
	if code != 0 || stdout != want {
		t.Errorf("status %d, output\n%s\nwant\n%s\n%s", code, stdout, want, stderr)
	}
}

// TestProgramsEndAsGoDoes runs programs in testdata/exits, which end
// badly, under go run and Deno: the same exit status, the same output,
// and the line of stderr that says why.
func TestProgramsEndAsGoDoes(t *testing.T) {
	t.Parallel()
	run := jsRuntime(t)
	for _, tc := range []struct{ dir, why string }{
		{"deadlock", "fatal error: all goroutines are asleep - deadlock!"},
		{"gopanic", "panic: boom 42"},
	} {
		t.Run(tc.dir, func(t *testing.T) {
			dir := "./testdata/exits/" + tc.dir
			goOut, goErr, goCode := runProgram(exec.Command("go", "run", dir))
			out := programs(t)
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
