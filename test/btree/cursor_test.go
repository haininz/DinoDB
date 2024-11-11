package btree_test

// func testBTreeCursorAt(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Insert entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Insert(i, i%btree_salt)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve entries
// 	for i := int64(0); i <= 1000; i++ {
// 		cur, err := index.CursorAt(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		entry, err := cur.GetEntry()
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }
