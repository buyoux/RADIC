package test

import (
	"RADIC/course"
	"fmt"
	"testing"
)

func TestBuildInvertIndex(t *testing.T) {
	docs := []*course.Doc{{Id: 1, Keywords: []string{"go", "数据库", "系统"}},
		{Id: 2, Keywords: []string{"go", "数据结构", "算法", "系统"}}}
	index := course.BuildInvertIndex(docs)
	for key, value := range index {
		fmt.Println("Keyword:", key, ", Id:", value)
	}
}

// go test -v .\course\test\ -run=^TestBuildInvertIndex$ -count=1
