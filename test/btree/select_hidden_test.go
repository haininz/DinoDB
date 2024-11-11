package btree_test

import (
	"dinodb/pkg/btree"
	"dinodb/test/utils"
	"fmt"
	"testing"
)

func TestBTreeSelectHidden(t *testing.T) {
	t.Run("Increasing/Thousand", stageSelectIncreasingTest(1000))
	t.Run("DeletedEntriesNotFound", testSelectDeletedEntriesNotFound)
}

func TestBTreeSelectRangeHidden(t *testing.T) {
	t.Run("All", testSelectRangeAll)
	t.Run("Empty", testSelectRangeEmpty)
	t.Run("EmptyNode", testSelectRangeEmptyNode)
}

/*
Creates a BTree index, inserts 1000 entries, deletes some entries,
and makes sure deleted entries are not found in FindRange
*/
func testSelectDeletedEntriesNotFound(t *testing.T) {
	initialNumEntries := int64(1000)
	index := standardBTreeSetup(t, initialNumEntries)

	// Removes entries 200 to 499
	amountToDelete := int64(300)
	for i := range amountToDelete {
		err := index.Delete(i + 200)
		if err != nil {
			t.Error(err)
		}
	}
	// Retrieve all entries using TableFindRange
	entries, err := index.Select()
	if err != nil {
		t.Error(err)
	}
	expectedLenEntries := (initialNumEntries - amountToDelete)
	//check that size of entries slice is expected
	if int64(len(entries)) != expectedLenEntries {
		err = fmt.Errorf("Wrong number of entries returned by TableFindRange; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), expectedLenEntries)
		t.Error(err)
	}
	//check that none of the entries are the deleted ones
	for _, entry := range entries {
		if entry.Key >= int64(200) && entry.Key < int64(500) {
			t.Error("Deleted entry found in slice returned from TableFindRange")
			break
		}
	}
	index.Close()
}

/*
Creates a BTree index, inserts 1000 entries, and then retrieves all the
entries through SelectRange
*/
func testSelectRangeAll(t *testing.T) {
	index := standardBTreeSetup(t, 1000)

	// Retrieve entries
	start := int64(0)
	end := int64(1000)
	entries, err := index.SelectRange(start, end)
	if err != nil {
		t.Error(err)
	}
	//check that size of entries slice is expected
	expectedLenEntries := (end - start)
	if int64(len(entries)) != expectedLenEntries {
		err = fmt.Errorf("Wrong number of entries returned by SelectRange; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), expectedLenEntries)
		t.Error(err)
	}
	for i, entry := range entries {
		key := int64(i)
		utils.CheckEntry(t, entry, key, generateValue(key))
	}
	index.Close()
}

/*
Tests that SelectRange returns no entries when called on an empty range
*/
func testSelectRangeEmpty(t *testing.T) {
	index := setupBTree(t)

	// Insert boundary entries
	utils.InsertEntry(t, index, 199, generateValue(199))
	utils.InsertEntry(t, index, 301, generateValue(301))

	// Call TableFindRange on an empty range
	// Call SelectRange on an empty ange
	start := int64(200)
	end := int64(300)
	entries, err := index.SelectRange(start, end)
	if err != nil {
		t.Error(err)
	}
	if len(entries) != 0 {
		t.Errorf("SelectRange called on an empty range returned entries. Expected number of entries: 0, actual: %d", len(entries))
	}

	index.Close()
}

/*
Creates a BTree index, inserts 1000 entries, deletes enough entries to make empty nodes,
and then retrieves all the entries through SelectRange
*/
func testSelectRangeEmptyNode(t *testing.T) {
	index := standardBTreeSetup(t, 1000)

	// Remove entries in a middle node
	// Removes all entries from Node #2 --- entries 101 inclusive to 202 exclusive
	for i := btree.ENTRIES_PER_LEAF_NODE / 2; i < btree.ENTRIES_PER_LEAF_NODE; i++ {
		err := index.Delete(i)
		if err != nil {
			t.Error(err)
		}
	}
	//Check that we can still retrieve all other entries contiguously
	// Retrieve entries
	start := int64(0)
	end := int64(1000)
	entries, err := index.SelectRange(start, end)
	if err != nil {
		t.Error(err)
	}
	//check that size of entries slice is expected
	expectedLenEntries := ((end - start) - (btree.ENTRIES_PER_LEAF_NODE - (btree.ENTRIES_PER_LEAF_NODE / 2)))
	if int64(len(entries)) != expectedLenEntries {
		err = fmt.Errorf("Wrong number of entries returned by SelectRange; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), expectedLenEntries)
		t.Error(err)
	}
	//check that the entries returned match expected entries
	for i := range btree.ENTRIES_PER_LEAF_NODE / 2 {
		entry := entries[i]
		key := int64(i)
		utils.CheckEntry(t, entry, key, generateValue(key))
	}
	for i := btree.ENTRIES_PER_LEAF_NODE / 2; i < (1000 - (btree.ENTRIES_PER_LEAF_NODE - (btree.ENTRIES_PER_LEAF_NODE / 2))); i++ {
		entry := entries[i]
		key := i + (btree.ENTRIES_PER_LEAF_NODE / 2)
		utils.CheckEntry(t, entry, key, generateValue(key))
	}
	index.Close()
}
