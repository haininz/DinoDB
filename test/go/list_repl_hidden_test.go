package go_test

import (
	"dinodb/pkg/list"
	"io"
	"testing"
)

func TestListReplHidden(t *testing.T) {
	t.Run("AddTail", testListReplAddTail)
	t.Run("IntegratedPushes", testListReplIntegratedPushes)
	t.Run("Remove", testListReplRemove)
}

// Tests successful and failed list_push_tail calls
func testListReplAddTail(t *testing.T) {
	inputWriter, output := startListRepl(t)

	//successful list_push_tail
	io.WriteString(inputWriter, "list_push_tail 1\n")
	checkSuccessOutput(t, output, "list_push_tail")

	//illformed list_push_tail
	io.WriteString(inputWriter, "list_push_tail\n")
	checkOutputHasErrorMessage(t, output, list.ErrListPushTailInvalidArgs)
}

// Tests that list_push_head and list_push_tail work together properly
func testListReplIntegratedPushes(t *testing.T) {
	inputWriter, output := startListRepl(t)

	//successful list_push_tail
	io.WriteString(inputWriter, "list_push_tail 1\n")
	checkSuccessOutput(t, output, "list_push_tail")

	//successful list_push_head
	io.WriteString(inputWriter, "list_push_head 2\n")
	checkSuccessOutput(t, output, "list_push_head")

	//successful list_push_tail
	io.WriteString(inputWriter, "list_push_tail 3\n")
	checkSuccessOutput(t, output, "list_push_tail")

	//print out and check order
	io.WriteString(inputWriter, "list_print\n")
	checkOutputExact(t, output, "2\n1\n3\n")
}

// Tests successful and failed list_remove calls
func testListReplRemove(t *testing.T) {
	inputWriter, output := startListRepl(t)

	//add an element to list to remove
	io.WriteString(inputWriter, "list_push_head 1\n")

	//ill-formed list_remove call
	io.WriteString(inputWriter, "list_remove\n")
	checkOutputHasErrorMessage(t, output, list.ErrListRemoveInvalidArgs)

	//successful list_remove call
	io.WriteString(inputWriter, "list_remove 1\n")
	checkSuccessOutput(t, output, "list_remove")

	//failed list_remove call on non-existent link
	io.WriteString(inputWriter, "list_remove 1\n")
	checkOutputHasErrorMessage(t, output, list.ErrListRemoveValueNotFound)
}
