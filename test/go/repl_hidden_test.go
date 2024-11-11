package go_test

import (
	"dinodb/pkg/repl"
	"fmt"
	"strings"
	"testing"
)

func TestReplHidden(t *testing.T) {
	t.Run("AddStress", testAddStress)
	t.Run("CombineOneRepl", testCombineOneRepl)
	t.Run("Combine", testCombine)
	t.Run("CombineDuplicateTrigger", testCombineDuplicateTrigger)
}

// Tests that commands and help strings are valid upon adding mass amounts of commands (x1000)
func testAddStress(t *testing.T) {
	r := repl.NewRepl()
	r.AddCommand("1", f1, "1 help")
	for i := 0; i < 1000; i++ {
		r.AddCommand(fmt.Sprintf("%v", i), f1, fmt.Sprintf("%v help string", i))
	}
	for i := 0; i < 1000; i++ {
		if _, ok := r.GetCommands()[fmt.Sprintf("%v", i)]; !ok {
			t.Fatal("bad add command")
		}
		if _, ok := r.GetHelp()[fmt.Sprintf("%v", i)]; !ok {
			t.Fatal("bad add command")
		}
	}
}

// Tests that CombineRepls() works properly when combining just one REPL with nothing else
func testCombineOneRepl(t *testing.T) {
	r1 := repl.NewRepl()
	r1.AddCommand("1", f1, "1 help")
	r1.AddCommand("2", f2, "2 help")
	r, err := repl.CombineRepls([]*repl.REPL{r1})
	if err != nil {
		t.Fatal("bad combine")
	}
	if _, ok := r.GetCommands()["1"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetCommands()["2"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetHelp()["1"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if _, ok := r.GetHelp()["2"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
}

/*
Tests that CombineRepl() works properly when combining 2 REPLs that
each have existing commands and help strings in them.
*/
func testCombine(t *testing.T) {
	r1 := repl.NewRepl()
	r1.AddCommand("1", f1, "1 help")
	r1.AddCommand("2", f2, "2 help")
	r2 := repl.NewRepl()
	r2.AddCommand("3", f3, "3 help")
	r2.AddCommand("4", f4, "4 help")
	r2.AddCommand("5", f5, "5 help")
	r, err := repl.CombineRepls([]*repl.REPL{r1, r2})
	if err != nil {
		t.Fatal("bad combine")
	}
	if _, ok := r.GetCommands()["1"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetCommands()["2"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetCommands()["3"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetCommands()["4"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetCommands()["5"]; !ok {
		t.Fatal("bad combine - funcs malformed")
	}
	if _, ok := r.GetHelp()["1"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if _, ok := r.GetHelp()["2"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if _, ok := r.GetHelp()["3"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if _, ok := r.GetHelp()["4"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if _, ok := r.GetHelp()["5"]; !ok {
		t.Fatal("bad combine - help malformed")
	}
	if !strings.Contains(r.HelpString(), "1 help") {
		t.Fatal("bad combine - print help malformed")
	}
	if !strings.Contains(r.HelpString(), "2 help") {
		t.Fatal("bad combine - print help malformed")
	}
	if !strings.Contains(r.HelpString(), "3 help") {
		t.Fatal("bad combine - print help malformed")
	}
	if !strings.Contains(r.HelpString(), "4 help") {
		t.Fatal("bad combine - print help malformed")
	}
	if !strings.Contains(r.HelpString(), "5 help") {
		t.Fatal("bad combine - print help malformed")
	}
}

/*
Tests that CombineRepl() works properly when combining 2 REPLs
that have overlapping commands and help strings in them by erroring
and not combining.
*/
func testCombineDuplicateTrigger(t *testing.T) {
	r1 := repl.NewRepl()
	r1.AddCommand("1", f1, "1 help")
	r2 := repl.NewRepl()
	r2.AddCommand("1", f1, "1 again")
	_, err := repl.CombineRepls([]*repl.REPL{r1, r2})
	if err == nil {
		t.Fatal("should error if duplicate triggers")
	}
}

func TestReplRunHidden(t *testing.T) {
	t.Run("OverwriteCommand", testRunOverwriteCommand)
	t.Run("MultipleCommands", testRunMultipleCommands)
	t.Run("MultipleCommandsHelp", testRunMultipleCommandsHelp)
}

func testRunOverwriteCommand(t *testing.T) {
	r := repl.NewRepl()
	r.AddCommand("echo", echo, "prints back everything")
	r.AddCommand("echo", f1, "f1 help")
	input, output := startRepl(t, r)

	fmt.Fprintln(input, "echo hey")
	checkOutputExact(t, output, "")
}

func testRunMultipleCommands(t *testing.T) {
	r := repl.NewRepl()
	r.AddCommand("hey", func(s string, r *repl.REPLConfig) (output string, err error) {
		return "hey", nil
	}, "says hey")
	r.AddCommand("echo", echo, "prints back everything")
	r.AddCommand("smush", func(s string, r *repl.REPLConfig) (output string, err error) {
		return strings.Join(strings.Fields(s), ""), nil
	}, "removes all whitespaces")
	input, output := startRepl(t, r)

	// try running all possible commands and checking output

	fmt.Fprintln(input, "hey")
	checkOutputExact(t, output, "hey\n")

	fmt.Fprintln(input, "echo based")
	checkOutputExact(t, output, "echo based\n")

	fmt.Fprintln(input, "smush hello\tworld !")
	checkOutputExact(t, output, "smushhelloworld!\n")
}

func testRunMultipleCommandsHelp(t *testing.T) {
	helpMap := map[string]string{
		"hey":   "says hey",
		"echo":  "prints back everything",
		"smush": "removes all whitespaces",
	}

	r := repl.NewRepl()
	r.AddCommand("hey", func(s string, r *repl.REPLConfig) (output string, err error) {
		return "hey", nil
	}, helpMap["hey"])
	r.AddCommand("echo", echo, helpMap["echo"])
	r.AddCommand("smush", func(s string, r *repl.REPLConfig) (output string, err error) {
		return strings.Join(strings.Fields(s), ""), nil
	}, helpMap["smush"])
	input, output := startRepl(t, r)

	checkHelp(t, input, output, helpMap)
}
