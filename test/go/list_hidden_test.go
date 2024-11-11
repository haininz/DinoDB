package go_test

import (
	"dinodb/pkg/list"
	"testing"
)

func TestListHidden(t *testing.T) {
	t.Run("SingletonListTail", testSingletonListTail)
	t.Run("PushHeadTailIntList", testPushHeadTailIntList)
	t.Run("MapEmptyListNoSideEffect", testMapEmptyListNoSideEffect)
	t.Run("PopNewTail", testPopNewTail)
	t.Run("PopSelfEmpty", testPopSelfEmpty)
	t.Run("PopSelfThree", testPopSelfThree)
	t.Run("PopSelfCleanup", testPopSelfCleanup)
}

/*
Also tests that in a list with only one element,
the head of the list is the same as the tail of the list,
but this time using PushTail instead of PushHead
*/
func testSingletonListTail(t *testing.T) {
	list := list.NewList()
	list.PushTail(5)
	if list.PeekHead() != list.PeekTail() {
		t.Fatal("head not equal to tail in singleton list")
	}
}

/*
Adds multiple elements to list by alternating pushing to head
and tail, and then tests that the order of elements is correct
and that the head and tail values are correct.
*/
func testPushHeadTailIntList(t *testing.T) {
	l := list.NewList()
	l.PushTail(1)
	l.PushHead(2)
	l.PushTail(3)
	l.PushHead(4)
	l.PushTail(5)
	if l.PeekHead() == nil || l.PeekHead().GetValue() != 4 {
		t.Fatal("bad peekhead")
	}
	if l.PeekTail() == nil || l.PeekTail().GetValue() != 5 {
		t.Fatal("bad peektail")
	}
	verifyList(t, l, []interface{}{4, 2, 1, 3, 5})
}

/*
Tests that Map() does not have any effect when there are no
links in the list it is called on.
*/
func testMapEmptyListNoSideEffect(t *testing.T) {
	l := list.NewList()
	x := 0
	lambda := func(_ *list.Link) { x += 1 }
	l.Map(lambda)
	if x != 0 {
		t.Fatal("map should have run 0 times on empty list")
	}
}

/*
Tests that the head and tail of the list update
properly when PopSelf() is called on the tail of a list.
*/
func testPopNewTail(t *testing.T) {
	l := list.NewList()
	l.PushHead(2)
	l.PushHead(1)
	elt1 := l.Find(func(x *list.Link) bool { return x.GetValue() == 1 })
	elt2 := l.Find(func(x *list.Link) bool { return x.GetValue() == 2 })
	elt2.PopSelf()
	if l.PeekHead() != elt1 {
		t.Fatal("bad pop, head not updated")
	}
	if l.PeekTail() != elt1 {
		t.Fatal("bad pop, tail not updated")
	}
}

/*
Tests that list and link metadata are properly updated when
PopSelf() is called on the only link in a list.
*/
func testPopSelfEmpty(t *testing.T) {
	l := list.NewList()
	l.PushHead(1)
	elt := l.PeekHead()
	elt.PopSelf()
	if elt.GetNext() != nil || elt.GetPrev() != nil {
		t.Fatal("bad pop; elt still has next or prev")
	}
	if l.PeekHead() != nil || l.PeekTail() != nil {
		t.Fatal("bad pop; list head or tail should be nil")
	}
}

/*
Tests that link metadata are properly updated when PopSelf()
is called on a link in the middle of a list.
*/
func testPopSelfThree(t *testing.T) {
	l := list.NewList()
	l.PushHead(1)
	l.PushHead(2)
	l.PushHead(3)
	elt1 := l.Find(func(x *list.Link) bool { return x.GetValue() == 1 })
	elt2 := l.Find(func(x *list.Link) bool { return x.GetValue() == 2 })
	elt3 := l.Find(func(x *list.Link) bool { return x.GetValue() == 3 })
	elt2.PopSelf()
	if elt1.GetPrev() != elt3 || elt3.GetNext() != elt1 {
		t.Fatal("bad pop; next or prev aren't pointed correctly")
	}
}

/* Tests that link metadata is cleaned up properly after PopSelf() is called. */
func testPopSelfCleanup(t *testing.T) {
	l := list.NewList()
	l.PushHead(1)
	l.PushHead(2)
	l.PushHead(3)
	elt2 := l.Find(func(x *list.Link) bool { return x.GetValue() == 2 })
	elt2.PopSelf()
	if elt2.GetList() != nil || elt2.GetNext() != nil || elt2.GetPrev() != nil {
		t.Fatal("popped link's metadata is not cleaned up properly")
	}
}

// func testDoublePopSelf(t *testing.T) {
// 	l := list.NewList()
// 	l.PushHead(1)
// 	elt := l.PeekHead()
// 	elt.PopSelf()
// 	elt.PopSelf()
// 	if elt.GetNext() != nil || elt.GetPrev() != nil {
// 		t.Fatal("bad double pop; should be idempotent")
// 	}
// 	if l.PeekHead() != nil || l.PeekTail() != nil {
// 		t.Fatal("bad double pop; should be idempotent")
// 	}
// }
