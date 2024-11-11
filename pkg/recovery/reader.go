package recovery

import (
	"bytes"
	"io"

	"github.com/google/uuid"
	"github.com/icza/backscanner"
)

// Helper method that gets all log strings and most recent checkpoint position from the log file.
func (rm *RecoveryManager) getRelevantStrings() (
	relevantStrings []string, checkpointPos int, err error) {
	fstats, err := rm.fd.Stat()
	if err != nil {
		return nil, 0, err
	}

	scanner := backscanner.New(rm.fd, int(fstats.Size()))
	checkpointTarget := []byte("checkpoint")
	startTarget := []byte("start")
	relevantStrings = make([]string, 0)
	checkpointHit := false
	txs := make(map[uuid.UUID]bool)
	for {
		line, _, err := scanner.LineBytes()
		if err != nil {
			if err == io.EOF {
				return relevantStrings, 0, nil
			} else {
				return nil, 0, err
			}
		}
		relevantStrings = append([]string{string(line)}, relevantStrings...)
		checkpointPos += 1
		if checkpointHit {
			if bytes.Contains(line, startTarget) {
				log, err := logFromString(string(line))
				if err != nil {
					return nil, 0, err
				}
				id := log.(startLog).id
				delete(txs, id)
			}
		}
		if !checkpointHit && bytes.Contains(line, checkpointTarget) {
			checkpointHit = true
			log, err := logFromString(string(line))
			if err != nil {
				return nil, 0, err
			}
			for _, tx := range log.(checkpointLog).ids {
				txs[tx] = true
			}
			checkpointPos = 0
		}
		if checkpointHit && len(txs) <= 0 {
			break
		}
	}
	return relevantStrings, checkpointPos, err
}

// Reads in ALL the logs and most recent checkpoint position from disk.
func (rm *RecoveryManager) readLogs() (logs []log, checkpointPos int, err error) {
	strings, checkpointPos, err := rm.getRelevantStrings()
	if err != nil {
		return nil, 0, err
	}
	if len(strings) > 0 {
		logs = make([]log, len(strings)-1)
		for i, s := range strings[:len(strings)-1] {
			log, err := logFromString(s)
			if err != nil {
				return nil, 0, err
			}
			logs[i] = log
		}
	} else {
		logs = make([]log, 0)
	}
	return logs, checkpointPos, nil
}
