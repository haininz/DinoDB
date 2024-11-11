package concurrency_test

import (
	"dinodb/pkg/concurrency"
	"testing"
)

func TestDeadlockHidden(t *testing.T) {
	t.Run("DAG", testDeadlockDAG)
	t.Run("Triangle", testDeadlockTriangle)
	t.Run("SelfLoop", testDeadlockSelfLoop)
	t.Run("Disjoint", testDeadlockDisjoint)
}

func testDeadlockDAG(t *testing.T) {
	t1 := concurrency.Transaction{}
	t2 := concurrency.Transaction{}
	t3 := concurrency.Transaction{}
	t4 := concurrency.Transaction{}
	g := concurrency.NewGraph()
	g.AddEdge(&t1, &t2)
	g.AddEdge(&t1, &t3)
	g.AddEdge(&t2, &t4)
	g.AddEdge(&t3, &t4)
	if g.DetectCycle() {
		t.Error("cycle detected in DAG")
	}
}

func testDeadlockTriangle(t *testing.T) {
	t1 := concurrency.Transaction{}
	t2 := concurrency.Transaction{}
	t3 := concurrency.Transaction{}
	g := concurrency.NewGraph()
	g.AddEdge(&t1, &t2)
	g.AddEdge(&t2, &t3)
	g.AddEdge(&t3, &t1)
	if !g.DetectCycle() {
		t.Error("could not detect triangle cycle")
	}
}

func testDeadlockSelfLoop(t *testing.T) {
	t1 := concurrency.Transaction{}
	g := concurrency.NewGraph()
	g.AddEdge(&t1, &t1)
	if !g.DetectCycle() {
		t.Error("could not detect self-loop")
	}
}

func testDeadlockDisjoint(t *testing.T) {
	t1 := concurrency.Transaction{}
	t2 := concurrency.Transaction{}
	t3 := concurrency.Transaction{}
	t4 := concurrency.Transaction{}
	g := concurrency.NewGraph()
	g.AddEdge(&t1, &t2)
	g.AddEdge(&t2, &t1)
	g.AddEdge(&t3, &t4)
	g.AddEdge(&t4, &t3)
	if !g.DetectCycle() {
		t.Error("failed to detect cycle")
	}
}
