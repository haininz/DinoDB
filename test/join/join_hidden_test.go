package join_test

import (
	"dinodb/pkg/entry"
	"dinodb/pkg/join"
	"testing"
)

func TestJoinHidden(t *testing.T) {
	t.Run("ValueOnKey", testJoinValueOnKey)
	t.Run("KeysManyToMany", testJoinKeysManyToMany)
	t.Run("KeysOneToMany", testJoinKeysOneToMany)
}

// join on left value and right key
func testJoinValueOnKey(t *testing.T) {
	index1, index2 := setupJoin(t)

	//insert right index entries
	insertIntoIndex(t, index2, 5, 5%join_salt)
	insertIntoIndex(t, index2, 6%join_salt, 6)

	// Check for expected values of results
	// Use a map as a set to check the existence of all values
	expectedResultsMap := make(map[string]struct{})
	// Insert left index entries and store expected results
	for i := int64(0); i < 10; i++ {
		insertIntoIndex(t, index1, i, i%join_salt)
		if (i % join_salt) == (6 % join_salt) {
			lResult := entry.New(i, i%join_salt)
			rResult := entry.New(6%join_salt, 6)
			entryPair := join.EntryPair{L: lResult, R: rResult}
			key := entryPairKey(entryPair)
			expectedResultsMap[key] = struct{}{}
		}
		if (i % join_salt) == 5 {
			lResult := entry.New(i, i%join_salt)
			rResult := entry.New(5, 5%join_salt)
			entryPair := join.EntryPair{L: lResult, R: rResult}
			key := entryPairKey(entryPair)
			expectedResultsMap[key] = struct{}{}
		}
	}

	// Get and check results.
	results := getResults(t, index1, index2, false, true)
	if len(results) != 2 {
		t.Errorf("basic join not working; expected %v results, got %d\n", 2, len(results))
	}

	resultsMap := make(map[string]struct{})
	for _, result := range results {
		key := entryPairKey(result)
		resultsMap[key] = struct{}{}
	}

	equal := equalEntryPairs(expectedResultsMap, resultsMap)
	if !equal {
		t.Errorf("join results have incorrect values")
	}
}

// stress test joining on keys with multiple matching pairs
// does not check values of the produced matches
func testJoinKeysManyToMany(t *testing.T) {
	tests := map[string]int64{
		"Ten":    10,
		"Stress": 128,
	}

	for name, numInserts := range tests {
		t.Run(name, func(t *testing.T) {
			index1, index2 := setupJoin(t)

			// Insert entries. Since we have duplicate keys, only BTreeIndex cannot be used
			for i := int64(0); i < numInserts; i++ {
				insertIntoIndex(t, index1, 0, i%join_salt)
				insertIntoIndex(t, index2, 0, i%join_salt)
			}

			// Get and check results.
			results := getResults(t, index1, index2, true, true)
			expectedNumEntries := numInserts * numInserts
			if int64(len(results)) != expectedNumEntries {
				t.Errorf("basic join not working; expected %v results, got %d\n", expectedNumEntries, len(results))
			}
		})
	}
}

// check that a join with multiple matching keys produces all valid resulting pairs
// also checks the value of the produced pairs
func testJoinKeysOneToMany(t *testing.T) {
	index1, index2 := setupJoin(t)

	//insert one entry into left index
	insertIntoIndex(t, index1, 0, 11)

	// Check for expected values of results
	// Use a map as a set to check the existence of all values
	expectedResultsMap := make(map[string]struct{})
	// Insert entries into right index.
	for i := int64(0); i < 10; i++ {
		insertIntoIndex(t, index2, 0, i%join_salt)
		lResult := entry.New(0, 11)
		rResult := entry.New(0, i%join_salt)
		entryPair := join.EntryPair{L: lResult, R: rResult}
		key := entryPairKey(entryPair)
		expectedResultsMap[key] = struct{}{}
	}

	// Get and check results.
	results := getResults(t, index1, index2, true, true)
	if len(results) != 10 {
		t.Errorf("basic join not working; expected %v results, got %d\n", 0, len(results))
	}

	//check results values
	resultsMap := make(map[string]struct{})
	for _, result := range results {
		key := entryPairKey(result)
		resultsMap[key] = struct{}{}
	}

	equal := equalEntryPairs(expectedResultsMap, resultsMap)
	if !equal {
		t.Errorf("join results have incorrect values")
	}
}
