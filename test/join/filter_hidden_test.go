package join_test

import (
	"dinodb/pkg/join"
	"testing"
)

func TestBloomFilterHidden(t *testing.T) {
	t.Run("Stress", testFilterStress)
	t.Run("InsertAndCheckBig", testFilterInsertAndCheckBig)
}

func testFilterInsertAndCheckBig(t *testing.T) {
	filter := join.CreateFilter(2048)
	for i := 0; i < 100; i++ {
		filter.Insert(int64(i))
		if !filter.Contains(int64(i)) {
			t.Errorf("inserted value %d but not found", i)
		}
	}
}

func testFilterStress(t *testing.T) {
	filter := join.CreateFilter(65536)
	for i := 0; i < 16384; i++ {
		filter.Insert(int64(i))
	}
	for i := 0; i < 16384; i++ {
		if !filter.Contains(int64(i)) {
			t.Errorf("should have 'found' value %d", i)
		}
	}
}
