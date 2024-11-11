package btree_test

import "testing"

func TestBTreeInsertHidden(t *testing.T) {
	t.Run("Ascending", testHiddenInsertAscending)
	t.Run("Random", testHiddenInsertRandom)
}

func testHiddenInsertAscending(t *testing.T) {
	// Define the test cases
	tests := map[string]InsertTestData{
		"StressNoWrite":   {1_000_000, false},
		"StressWithWrite": {1_000_000, true},
	}

	// Runs the test cases
	for name, testData := range tests {
		t.Run(name, stageInsertAscending(testData))
	}
}

func testHiddenInsertRandom(t *testing.T) {
	// Define the test cases
	tests := map[string]InsertTestData{
		"StressNoWrite":   {20_000, false},
		"StressWithWrite": {20_000, true},
	}

	// Run the test cases
	for name, testData := range tests {
		t.Run(name, stageInsertRandom(testData))
	}
}
