package test

import (
	"RADIC/internal/kvdb"
	"RADIC/util"
	"testing"
)

func TestBadger(t *testing.T) {
	setup = func() {
		var err error
		db, err = kvdb.GetKvDb(kvdb.BADGER, util.ProjectRootPath+"data/badger_db")
		if err != nil {
			panic(err)
		}
	}
	t.Run("badger test", testPipeline)
}

// go test -v .\internal\kvdb\test\ -run=^TestBadger$ -count=1
