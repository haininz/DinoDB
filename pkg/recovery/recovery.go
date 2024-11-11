package recovery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"dinodb/pkg/concurrency"
	"dinodb/pkg/config"
	"dinodb/pkg/database"

	"github.com/otiai10/copy"

	"github.com/google/uuid"
)

// Recovery Manager.
type RecoveryManager struct {
	db      *database.Database
	tm      *concurrency.TransactionManager
	txStack map[uuid.UUID][]log
	fd      *os.File
	mtx     sync.Mutex
}

// Construct a recovery manager.
func NewRecoveryManager(
	db *database.Database,
	tm *concurrency.TransactionManager,
	logName string,
) (rm *RecoveryManager, err error) {
	fd, err := os.OpenFile(logName, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &RecoveryManager{
		db:      db,
		tm:      tm,
		txStack: make(map[uuid.UUID][]log),
		fd:      fd,
	}, nil
}

// Write the string `s` to the log file. Expects rm.mtx to be locked
func (rm *RecoveryManager) writeToBuffer(s string) error {
	_, err := rm.fd.WriteString(s)
	if err != nil {
		return err
	}
	err = rm.fd.Sync()
	return err
}

// Write a Table log.
func (rm *RecoveryManager) Table(tblType string, tblName string) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	/* SOLUTION {{{ */
	tl := tableLog{
		tblType: tblType,
		tblName: tblName,
	}
	rm.writeToBuffer(tl.toString())
	/* SOLUTION }}} */
}

// Write an Edit log.
func (rm *RecoveryManager) Edit(clientId uuid.UUID, table database.Index, action action, key int64, oldval int64, newval int64) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	/* SOLUTION {{{ */
	el := editLog{
		id:        clientId,
		tablename: table.GetName(),
		action:    action,
		key:       key,
		oldval:    oldval,
		newval:    newval,
	}
	rm.txStack[clientId] = append(rm.txStack[clientId], el)
	rm.writeToBuffer(el.toString())
	/* SOLUTION }}} */
}

// Write a transaction start log.
func (rm *RecoveryManager) Start(clientId uuid.UUID) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	/* SOLUTION {{{ */
	sl := startLog{
		id: clientId,
	}
	rm.txStack[clientId] = make([]log, 0)
	rm.txStack[clientId] = append(rm.txStack[clientId], sl)
	rm.writeToBuffer(sl.toString())
	/* SOLUTION }}} */
}

// Write a transaction commit log.
func (rm *RecoveryManager) Commit(clientId uuid.UUID) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	/* SOLUTION {{{ */
	cl := commitLog{
		id: clientId,
	}
	delete(rm.txStack, clientId)
	rm.writeToBuffer(cl.toString())
	/* SOLUTION }}} */
}

// Flush all pages to disk and write a checkpoint log.
func (rm *RecoveryManager) Checkpoint() {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	/* SOLUTION {{{ */
	// Get a list of all running transactions.
	ids := make([]uuid.UUID, 0)
	// Shouldn't rm.tm.mtx be locked here?
	for _, tx := range rm.tm.GetTransactions() {
		ids = append(ids, tx.GetClientID())
	}
	// Create a checkpoint log.
	cl := checkpointLog{
		ids: ids,
	}
	tables := rm.db.GetTables()
	// Prevent all tables from being written to while checkpointing
	for _, v := range tables {
		v.GetPager().LockAllPages()
		defer v.GetPager().UnlockAllPages()
	}
	// Flush and write log.
	for _, v := range tables {
		v.GetPager().FlushAllPages()
	}
	rm.writeToBuffer(cl.toString())
	/* SOLUTION }}} */
	rm.delta() // Sorta-semi-pseudo-copy-on-write (to ensure db recoverability)
}

// Redo a given log's action.
func (rm *RecoveryManager) Redo(log log) error {
	switch log := log.(type) {
	case tableLog:
		payload := fmt.Sprintf("create %s table %s", log.tblType, log.tblName)
		_, err := database.HandleCreateTable(rm.db, payload)
		if err != nil {
			return err
		}
	case editLog:
		switch log.action {
		case INSERT_ACTION:
			payload := fmt.Sprintf("insert %v %v into %s", log.key, log.newval, log.tablename)
			err := database.HandleInsert(rm.db, payload)
			if err != nil {
				// There is already an entry, try updating
				payload := fmt.Sprintf("update %s %v %v", log.tablename, log.key, log.newval)
				err = database.HandleUpdate(rm.db, payload)
				if err != nil {
					return err
				}
			}
		case UPDATE_ACTION:
			payload := fmt.Sprintf("update %s %v %v", log.tablename, log.key, log.newval)
			err := database.HandleUpdate(rm.db, payload)
			if err != nil {
				// Entry may have been deleted, try inserting
				payload := fmt.Sprintf("insert %v %v into %s", log.key, log.newval, log.tablename)
				err := database.HandleInsert(rm.db, payload)
				if err != nil {
					return err
				}
			}
		case DELETE_ACTION:
			payload := fmt.Sprintf("delete %v from %s", log.key, log.tablename)
			err := database.HandleDelete(rm.db, payload)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("can only redo edit or table logs")
	}
	return nil
}

// Undo a given log's action.
func (rm *RecoveryManager) Undo(log editLog) error {
	switch log.action {
	case INSERT_ACTION:
		payload := fmt.Sprintf("delete %v from %s", log.key, log.tablename)
		err := HandleDelete(rm.db, rm.tm, rm, payload, log.id)
		if err != nil {
			return err
		}
	case UPDATE_ACTION:
		payload := fmt.Sprintf("update %s %v %v", log.tablename, log.key, log.oldval)
		err := HandleUpdate(rm.db, rm.tm, rm, payload, log.id)
		if err != nil {
			return err
		}
	case DELETE_ACTION:
		payload := fmt.Sprintf("insert %v %v into %s", log.key, log.oldval, log.tablename)
		err := HandleInsert(rm.db, rm.tm, rm, payload, log.id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Do a full recovery to the most recent checkpoint on startup.
func (rm *RecoveryManager) Recover() error {
	/* SOLUTION {{{ */
	logs, checkpointPos, err := rm.readLogs()
	if err != nil {
		return err
	}
	// Should exit if no logs/invalid checkpoint.
	if len(logs) <= checkpointPos {
		return nil
	}
	// Get uncommitted transactions at the checkpoint.
	activeTxs := make(map[uuid.UUID]bool)
	switch log := logs[checkpointPos].(type) {
	case checkpointLog:
		for _, id := range log.ids {
			activeTxs[id] = true
			rm.tm.Begin(id)
		}
	default:
	}
	// Redo everything after the checkpoint.
	logPtr := checkpointPos
	for ; logPtr < len(logs); logPtr++ {
		switch log := logs[logPtr].(type) {
		case tableLog:
			err := rm.Redo(log)
			if err != nil {
				return err
			}
		case editLog:
			err := rm.Redo(log)
			if err != nil {
				return err
			}
		case startLog:
			activeTxs[log.id] = true
			rm.tm.Begin(log.id)
		case commitLog:
			delete(activeTxs, log.id)
			rm.tm.Commit(log.id)
		default:
			continue
		}
	}
	// Undo uncommitted transactions.
	for logPtr = len(logs) - 1; logPtr >= 0; logPtr-- {
		if len(activeTxs) <= 0 {
			break
		}
		switch log := logs[logPtr].(type) {
		case editLog:
			if _, ok := activeTxs[log.id]; !ok {
				continue
			}
			err := rm.Undo(log)
			if err != nil {
				return err
			}
		case startLog:
			if _, ok := activeTxs[log.id]; !ok {
				continue
			}
			delete(activeTxs, log.id)
			rm.Commit(log.id)
			rm.tm.Commit(log.id)
		default:
			continue
		}
	}
	return nil
	/* SOLUTION }}} */
}

// Roll back the current uncommitted transaction for a client.
// This is called when you abort a transaction.
func (rm *RecoveryManager) Rollback(clientId uuid.UUID) error {
	/* SOLUTION {{{ */
	// Unwind the transaction actions stack and commit.
	logs := rm.txStack[clientId]
	// If no logs, immediately commit.
	if len(logs) == 0 {
		rm.Commit(clientId)
		rm.tm.Commit(clientId)
		return nil
	}
	// Check that our transaction stack is valid.
	switch logs[0].(type) {
	case startLog:
		break
	default:
		return errors.New("transaction stack must begin with a start log")
	}
	// Rollback the rest of the transactions FILO.
	for i := len(logs) - 1; i > 0; i-- {
		editLog, ok := logs[i].(editLog)
		if !ok {
			return errors.New("cannot undo non-edit log")
		}
		err := rm.Undo(editLog)
		if err != nil {
			return err
		}
	}
	// Write a commit, and release all the locks.
	rm.Commit(clientId)
	rm.tm.Commit(clientId)
	return nil
	/* SOLUTION }}} */
}

// Primes the database for recovery
func Prime(folder string) (*database.Database, error) {
	// Ensure folder is of the form */
	base := filepath.Clean(folder)
	recoveryFolder := base + "-recovery/"
	dbFolder := base + "/"

	// If recovery folder doesn't exist, create it and open db folder as normal
	if _, err := os.Stat(recoveryFolder); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(recoveryFolder, 0775)
			if err != nil {
				return nil, err
			}
			return database.Open(dbFolder)
		}
		return nil, err
	}

	// If recovery folder exists, replace db folder with recovery folder.
	// Copies over log file if it is in the db folder
	logSrcPath := filepath.Join(base, config.LogFileName)
	if _, err := os.Stat(logSrcPath); err == nil {
		logDstPath := filepath.Join(recoveryFolder, config.LogFileName)
		copy.Copy(logSrcPath, logDstPath)
	}
	os.RemoveAll(dbFolder)
	err := copy.Copy(recoveryFolder, dbFolder)
	if err != nil {
		return nil, err
	}
	return database.Open(dbFolder)
}

// Should be called at end of Checkpoint.
func (rm *RecoveryManager) delta() error {
	folder := strings.TrimSuffix(rm.db.GetBasePath(), "/")
	recoveryFolder := folder + "-recovery/"
	folder += "/"
	os.RemoveAll(recoveryFolder)
	err := copy.Copy(folder, recoveryFolder)
	return err
}
