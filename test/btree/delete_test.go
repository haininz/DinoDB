package btree_test

// func testBTreeDeleteTenNoWrite(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Insert entries
// 	for i := int64(0); i <= 10; i++ {
// 		err = index.Insert(i, i%btree_salt)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve entries
// 	for i := int64(0); i <= 10; i++ {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 		// Delete this entry
// 		err = index.Delete(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for i := int64(0); i <= 10; i++ {
// 		entry, err := index.Find(i)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteThousandNoWrite(t *testing.T) {
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
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 		// Delete this entry
// 		err = index.Delete(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteRandomNoWrite(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	nKeys := 1000

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Generate entries
// 	entries, answerKey := genRandomBTreeEntries(nKeys)
// 	for _, entry := range entries {
// 		key := entry.key
// 		val := entry.val
// 		err = index.Insert(key, val)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Insert duplicates
// 	for k := range answerKey {
// 		err = index.Insert(k, 0)
// 		if err == nil {
// 			t.Error("Could insert duplicate")
// 		}
// 	}
// 	// Delete random entries
// 	deleted := make(map[int64]int64)
// 	for k, v := range answerKey {
// 		shouldDelete := rand.Intn(2)
// 		if shouldDelete > 0 {
// 			err = index.Delete(k)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			deleted[k] = v
// 			delete(answerKey, k)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for k := range deleted {
// 		entry, err := index.Find(k)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	// Insert duplicates again
// 	for k := range answerKey {
// 		err = index.Insert(k, 0)
// 		if err == nil {
// 			t.Error("Could insert duplicate")
// 		}
// 	}
// 	// Retrieve entries
// 	for k, v := range answerKey {
// 		entry, err := index.Find(k)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != k {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != v {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for k := range deleted {
// 		entry, err := index.Find(k)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	// Insert deleted entries
// 	for k, v := range deleted {
// 		err = index.Insert(k, v)
// 		if err != nil {
// 			t.Error("Could not insert deleted entry")
// 		}
// 		answerKey[k] = v
// 	}
// 	// Retrieve entries
// 	for k, v := range answerKey {
// 		entry, err := index.Find(k)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != k {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != v {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteTen(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Insert entries
// 	for i := int64(0); i <= 10; i++ {
// 		err = index.Insert(i, i%btree_salt)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Close and reopen the database
// 	index.Close()
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Retrieve entries
// 	for i := int64(0); i <= 10; i++ {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 		// Delete this entry
// 		err = index.Delete(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for i := int64(0); i <= 10; i++ {
// 		entry, err := index.Find(i)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteThousand(t *testing.T) {
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
// 	//Find all 1000
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	// Close and reopen the database
// 	err = index.Close()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Retrieve entries and delete them
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 		// Delete this entry
// 		err = index.Delete(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteThousandReverse(t *testing.T) {
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
// 	//Find all 1000
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	// Close and reopen the database
// 	err = index.Close()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Retrieve entries and delete them in reverse order
// 	for i := int64(1000); i >= 0; i-- {
// 		entry, err := index.Find(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != i {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != i%btree_salt {
// 			t.Error("Entry found has the wrong value")
// 		}
// 		// Delete this entry
// 		err = index.Delete(i)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for i := int64(0); i <= 1000; i++ {
// 		entry, err := index.Find(i)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeDeleteRandom(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	nKeys := 1000

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Generate entries
// 	entries, answerKey := genRandomBTreeEntries(nKeys)
// 	for _, entry := range entries {
// 		key := entry.key
// 		val := entry.val
// 		err = index.Insert(key, val)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Insert duplicates
// 	for k := range answerKey {
// 		err = index.Insert(k, 0)
// 		if err == nil {
// 			t.Error("Could insert duplicate")
// 		}
// 	}
// 	// Delete random entries
// 	deleted := make(map[int64]int64)
// 	for k, v := range answerKey {
// 		shouldDelete := rand.Intn(2)
// 		if shouldDelete > 0 {
// 			err = index.Delete(k)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			deleted[k] = v
// 			delete(answerKey, k)
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for k := range deleted {
// 		entry, err := index.Find(k)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	// Close and reopen the database
// 	index.Close()
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Insert duplicates again
// 	for k := range answerKey {
// 		err = index.Insert(k, 0)
// 		if err == nil {
// 			t.Error("Could insert duplicate")
// 		}
// 	}
// 	// Retrieve entries
// 	for k, v := range answerKey {
// 		entry, err := index.Find(k)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != k {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != v {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	// Retrieve deleted entries
// 	for k := range deleted {
// 		entry, err := index.Find(k)
// 		if entry != nil || err == nil {
// 			t.Error("Could find deleted entry")
// 		}
// 	}
// 	// Insert deleted entries
// 	for k, v := range deleted {
// 		err = index.Insert(k, v)
// 		if err != nil {
// 			t.Error("Could not insert deleted entry")
// 		}
// 		answerKey[k] = v
// 	}
// 	// Retrieve entries
// 	for k, v := range answerKey {
// 		entry, err := index.Find(k)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if entry == nil {
// 			t.Error("Inserted entry could not be found")
// 		}
// 		if entry.Key != k {
// 			t.Error("Entry with wrong entry was found")
// 		}
// 		if entry.Value != v {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }
