package test

import (
	"RADIC/internal/kvdb"
	"RADIC/util"
	"testing"
)

func TestBolt(t *testing.T) {
	setup = func() {
		var err error
		// 使用工厂模式初始化bolt kvdb
		db, err = kvdb.GetKvDb(kvdb.BOLT, util.ProjectRootPath+"data/bolt_db")
		if err != nil {
			panic(err)
		}
	}
	t.Run("bolt_test", testPipeline)
}

// go test -v .\internal\kvdb\test\ -run=^TestBolt$ -count=1
