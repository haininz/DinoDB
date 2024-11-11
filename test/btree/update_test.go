package btree_test

// func testBTreeUpdateTenNoWrite(t *testing.T) {
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
// 	}
// 	// Update entries
// 	for i := int64(0); i <= 10; i++ {
// 		err = index.Update(i, -(i % btree_salt))
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateThousandNoWrite(t *testing.T) {
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
// 	// Update entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Update(i, -(i % btree_salt))
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateNonexistentNoWrite(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Update non-existent entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	// Insert entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Insert(i, i%btree_salt)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Update non-existent entries
// 	for i := int64(1001); i <= 2000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	// Update non-existent entries
// 	for i := int64(1001); i <= 2000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateRandomNoWrite(t *testing.T) {
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
// 	// Update entries
// 	for k := range answerKey {
// 		val := rand.Int63()
// 		err = index.Update(k, val)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		answerKey[k] = val
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

// func testBTreeUpdateTen(t *testing.T) {
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
// 	}
// 	// Update entries
// 	for i := int64(0); i <= 10; i++ {
// 		err = index.Update(i, -(i % btree_salt))
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateThousand(t *testing.T) {
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
// 	// Update entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Update(i, -(i % btree_salt))
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
// 		if entry.Value != -(i % btree_salt) {
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
// 		if entry.Value != -(i % btree_salt) {
// 			t.Error("Entry found has the wrong value")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateNonexistent(t *testing.T) {
// 	dbName := getTempDB(t)
// 	defer os.Remove(dbName)

// 	// Init the database
// 	index, err := btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Update non-existent entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	// Insert entries
// 	for i := int64(0); i <= 1000; i++ {
// 		err = index.Insert(i, i%btree_salt)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	// Update non-existent entries
// 	for i := int64(1001); i <= 2000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	// Close and reopen the database
// 	index.Close()
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Update non-existent entries
// 	for i := int64(1001); i <= 2000; i++ {
// 		err = index.Update(i, i%btree_salt)
// 		if err == nil {
// 			t.Error("Could update non-existent entry")
// 		}
// 	}
// 	index.Close()
// }

// func testBTreeUpdateRandom(t *testing.T) {
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
// 	// Update entries
// 	for k := range answerKey {
// 		val := rand.Int63()
// 		err = index.Update(k, val)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		answerKey[k] = val
// 	}
// 	// Close and reopen the database
// 	index.Close()
// 	index, err = btree.OpenTable(dbName)
// 	if err != nil {
// 		t.Error(err)
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
