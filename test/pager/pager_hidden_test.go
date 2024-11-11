package pager_test

import (
	"bytes"
	"dinodb/pkg/config"
	"testing"
)

func TestPagerHidden( t *testing.T) {
	// 6/16 Hidden Tests (38%)
	t.Run("GetDifferentPageData", testGetDifferentPageData)
	t.Run("GetPagesStress", testGetPagesStress)
	t.Run("FlushOnePageMultipleTimes", testFlushOnePageMultipleTimes)
	t.Run("GetExistingPage", testGetExistingPage)
	t.Run("FlushAllPages", testFlushAllPages)
	t.Run("GetNewPagesAfterPut", testGetNewPagesAfterPut)

}

/*
Writes to a new page without flushing.
Then ensures that future calls to GetNewPage returns a new page
WITHOUT the previously written data.
*/
func testGetDifferentPageData(t *testing.T) {
	p := setupPager(t)
	//get a page and write to it, but don't flush it
	p1 := getNewPage(t, p, true)
	data := []byte("test data")
	p1.Update(data, 0, int64(len(data)))
	//get a different page and check that the data is not the same
	p2 := getNewPage(t, p, true)
	// the data should NOT be the same
	if p1 == p2 {
		t.Error("Pages returned should not be the same")
	}
	if bytes.Equal(p2.GetData()[:len(data)], data) {
		t.Error("Data was written to other pages as well")
	}
}

/* Runs GetPage 10,000 times on the same page. */
func testGetPagesStress(t *testing.T) {
	p := setupPager(t)
	// Make a single new page.
	p1 := getNewPage(t, p, true)
	// Get the same page 10,000 times.
	for i := 0; i < 10000; i++ {
		page := getPage(t, p, p1.GetPageNum(), false)
		_ = p.PutPage(page)
	}
}

/*
Get a page, write into it, flush the data, and close the Pager.
Upon reopening, check that the data updated properly. Repeat this
process and get + write to the same page.
*/
func testFlushOnePageMultipleTimes(t *testing.T) {
	p := setupPager(t)
	// Write some data to page 0
	page := getNewPage(t, p, false)
	data := []byte("1")
	page.Update(data, 0, int64(len(data)))
	_ = p.PutPage(page)

	p.FlushPage(page)
	closeAndReopen(t, p)

	// Retrieve page and check that data is the same
	page = getPage(t, p, 0, false)
	if !bytes.Equal(page.GetData()[:len(data)], data) {
		t.Fatal("Initial update not flushed properly")
	}

	//write data at an offset after the previous data
	data = []byte("2")
	page.Update(data, 1, int64(len(data)))
	_ = p.PutPage(page)

	p.FlushPage(page)
	closeAndReopen(t, p)

	// Retrieve page and check that data from both updates are present
	page = getPage(t, p, 0, true)
	data = []byte("12")
	if !bytes.Equal(page.GetData()[:2], data) {
		t.Fatal("Second update was not flushed properly")
	}
}

/*
Checks that subsequent GetPage calls return the same empty page
as is returned by the initial GetNewPage call.
*/
func testGetExistingPage(t *testing.T) {
	p := setupPager(t)
	p1 := getNewPage(t, p, true)
	p2 := getPage(t, p, 0, true)
	p3 := getPage(t, p, 0, true)
	if p1 != p2 || p2 != p3 {
		t.Fatal("Should have retrieved page from the buffer cache; instead, page was duped")
	}
}

/*
Tests FlushAllPages by writing data to a variable number of new pages.
Then, after calling FlushAllPages and closing/reopening the pager, checks that the
same data previously written is still there.
*/
func testFlushAllPages(t *testing.T) {
	// Define the test cases. Maps subtest name to the number of pages to write to and flush
	tests := map[string]int{
		"MaxBufferPages": config.MaxPagesInBuffer,
		"Thousand":       1_000,
		"Stress":         10_000,
	}

	// Run the same test for each test case
	for name, numPages := range tests {
		t.Run(name, func(t *testing.T) {
			p := setupPager(t)
			data := []byte("hello")
			// Write some data to all pages
			for i := 0; i <= numPages; i++ {
				page := getNewPage(t, p, false)
				page.Update(data, 0, int64(len(data)))
				_ = p.PutPage(page)
			}

			p.FlushAllPages()
			closeAndReopen(t, p)

			// Check data is still in all pages
			for i := 0; i <= numPages; i++ {
				page := getPage(t, p, int64(i), false)
				if !bytes.Equal(page.GetData()[:len(data)], data) {
					t.Fatal("Data not flushed properly")
				}
				_ = p.PutPage(page)
			}
		})
	}
}

/*
Checks well-formedness of GetNewPage in relation to buffer cache size.
Fills up the active pages in the cache, then puts pages and sees if you 
can GetNewPage again.

Uses GetNewPage to get all the possible number of pages up to config.MaxPagesInBuffer
and checks that it works. Try to GetNewPage again and check that it fails and returns 
an error. Then, call PutPage followed by GetNewPage to check that higher pagenums
can be allocated past MaxPagesInBuffer.
*/
func testGetNewPagesAfterPut(t *testing.T) {
	p := setupPager(t)
	for i := 0; i < config.MaxPagesInBuffer - 1; i++ {
		_ = getNewPage(t, p, true)
	}
	prevPage := getNewPage(t, p, false)

	// Errors when GetNewPage() is called on full cache
	errPage, err := p.GetNewPage()
	if err == nil {
		_ = p.PutPage(errPage)
		_ = p.PutPage(prevPage)
		t.Fatal("Should have returned an error for running out of pages")
	}
	// Correctly puts the prev page to free up cache space
	err = p.PutPage(prevPage)
	if err != nil {
		t.Fatal("Should have allowed PutPage after GetNewPage errors")
	}

	// GetNewPage after a Page has been freed up
	prevPage, err = p.GetNewPage()
	if err != nil {
		t.Fatal("Should have allowed GetNewPage after space is freed by PutPage")
	}
	_ = p.PutPage(prevPage)
}