package hash_test

import (
	"dinodb/pkg/hash"
	"dinodb/test/utils"
	"testing"
)

func TestHashInsertHidden(t *testing.T) {
	// 5/12 Hidden tests (41%)
	t.Run("Ascending", testHiddenInsertAscending)
	t.Run("Random", testHiddenInsertRandom)
	t.Run("BucketSplit", testBucketSplit)
}

// Inserts a variable number of ascending keys and somewhat ascending values into a HashIndex,
// checking that they can be found with and without closing/flushing the index's data to disk
func testHiddenInsertAscending(t *testing.T) {
	// Define the test cases.
	insertAscendingTests := InsertTestsMap{
		"StressNoWrite":   {20_000, false},
		"StressWithWrite": {20_000, true},
	}

	// Run the test cases
	for name, testData := range insertAscendingTests {
		t.Run(name, stageInsertAscending(testData))
	}
}

func testHiddenInsertRandom(t *testing.T) {
	// Define the test cases.
	tests := InsertTestsMap{
		"StressNoWrite":   {20_000, false},
		"StressWithWrite": {20_000, true},
	}

	// Run the test cases
	for name, testData := range tests {
		t.Run(name, stageInsertRandom(testData))
	}
}

// Tests that buckets have less than MAX_BUCKET_SIZE entries,
// testing that buckets actually split.
func testBucketSplit(t *testing.T) {
	numInserts := int64(10_000)
	maxBucketSize := int(hash.MAX_BUCKET_SIZE)
	index := setupHash(t)

	// Insert entries
	for i := range numInserts {
		utils.InsertEntry(t, index, i, i*3)
	}

	pger := index.GetPager()
	// Check that each bucket contains less than 203
	bucketPns := index.GetTable().GetBuckets()
	for _, pn := range bucketPns {
		bucket, err := index.GetTable().GetBucketByPN(pn)
		if err != nil {
			t.Fatalf("Failed to retrieve bucket with pagenumber %d: %s", pn, err)
		}

		entries, err := bucket.Select()
		if err != nil {
			t.Fatal("Failed selecting all entries in a bucket:", err)
		}

		if len(entries) > maxBucketSize {
			t.Fatalf("Bucket had %d elements, should have split", len(entries))
		}

		// Release the pin on the bucket's page
		err = pger.PutPage(bucket.GetPage())
		if err != nil {
			t.Fatal("Failed to put the bucket's page:", err)
		}
	}
}
