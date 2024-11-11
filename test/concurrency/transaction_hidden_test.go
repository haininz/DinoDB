package concurrency_test

import (
	"dinodb/pkg/concurrency"
	"testing"
)

func TestTransactionHidden(t *testing.T) {
	t.Run("Advanced", testTransactionAdvanced)
	t.Run("LotsaLocks", testTransactionLotsaLocks)
	t.Run("ThreeProcessCycle", testTransactionThreeProcessCycle)
	t.Run("LongDAGNoCycle", testTransactionLongDAGNoCycle)
	t.Run("ReadLockNoCycleThreeProcess", testTransactionReadLockNoCycleThreeProcess)
	t.Run("LockIdentityAndIdempotency", testTransactionLockIdentityAndIdempotency)
}

func testTransactionAdvanced(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	// Sending instructions
	sendWithDelay(ch1, LockCommand{key: 0, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 2, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 3, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 0, lock: false, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: false, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 2, lock: false, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 3, lock: false, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{done: true})
	// Check for errors
	checkNoErrors(t, errch)
}

func testTransactionLotsaLocks(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	// Sending instructions
	for i := int64(0); i < 100; i++ {
		sendWithDelay(ch1, LockCommand{key: i, lock: true, lt: concurrency.W_LOCK})
	}
	sendWithDelay(ch1, LockCommand{done: true})
	// Check for errors
	checkNoErrors(t, errch)
}

func testTransactionThreeProcessCycle(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	tid2, ch2 := getTransactionThread()
	go handleTransactionThread(tm, index, tid2, ch2, errch)
	tid3, ch3 := getTransactionThread()
	go handleTransactionThread(tm, index, tid3, ch3, errch)
	// Sending instructions
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch2, LockCommand{key: 2, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch3, LockCommand{key: 3, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 2, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch2, LockCommand{key: 3, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch3, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{done: true})
	sendWithDelay(ch2, LockCommand{done: true})
	sendWithDelay(ch3, LockCommand{done: true})
	// Check for errors
	checkWasErrors(t, errch)
}

func testTransactionLongDAGNoCycle(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	tid2, ch2 := getTransactionThread()
	go handleTransactionThread(tm, index, tid2, ch2, errch)
	tid3, ch3 := getTransactionThread()
	go handleTransactionThread(tm, index, tid3, ch3, errch)
	tid4, ch4 := getTransactionThread()
	go handleTransactionThread(tm, index, tid4, ch4, errch)
	// Sending instructions
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch2, LockCommand{key: 2, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch3, LockCommand{key: 3, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch4, LockCommand{key: 4, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch2, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch3, LockCommand{key: 2, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch4, LockCommand{key: 3, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{done: true})
	sendWithDelay(ch2, LockCommand{done: true})
	sendWithDelay(ch3, LockCommand{done: true})
	sendWithDelay(ch4, LockCommand{done: true})
	// Check for errors
	checkNoErrors(t, errch)
}

func testTransactionReadLockNoCycleThreeProcess(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	tid2, ch2 := getTransactionThread()
	go handleTransactionThread(tm, index, tid2, ch2, errch)
	tid3, ch3 := getTransactionThread()
	go handleTransactionThread(tm, index, tid3, ch3, errch)
	// Sending instructions
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch2, LockCommand{key: 2, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch3, LockCommand{key: 3, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 2, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch2, LockCommand{key: 3, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch3, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{done: true})
	sendWithDelay(ch2, LockCommand{done: true})
	sendWithDelay(ch3, LockCommand{done: true})
	// Check for errors
	checkNoErrors(t, errch)
}

func testTransactionLockIdentityAndIdempotency(t *testing.T) {
	tm, index := setupTransaction(t)
	errch := make(chan error, BUFFER_SIZE)
	// Set up transactions
	tid1, ch1 := getTransactionThread()
	go handleTransactionThread(tm, index, tid1, ch1, errch)
	// Sending instructions
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.W_LOCK})
	sendWithDelay(ch1, LockCommand{key: 1, lock: true, lt: concurrency.R_LOCK})
	sendWithDelay(ch1, LockCommand{done: true})
	// Check for errors
	checkNoErrors(t, errch)
}
