package test

import (
	"fmt"
	"testing"

	"RADIC/course/dao"
	"RADIC/util"
)

var csvFile = util.ProjectRootPath + "data/bili_video.csv"

func TestDumpDataFromFile2DB1(t *testing.T) {
	dao.DumpDataFromFile2DB1(csvFile) //DumpDataFromFile2DB1 use time 117240 ms
	/*
		select count(*) from bili_video;
		delete from bili_video;
	*/
}

func TestDumpDataFromFile2DB2(t *testing.T) {
	dao.DumpDataFromFile2DB2(csvFile) //DumpDataFromFile2DB2 use time 7955 ms
	/*
		select count(*) from bili_video;
		delete from bili_video;
	*/
}

func TestDumpDataFromFile2DB3(t *testing.T) {
	dao.DumpDataFromFile2DB3(csvFile) //DumpDataFromFile2DB3 use time 3367 ms
	/*
		select count(*) from bili_video;
		delete from bili_video;
	*/
}

func testReadAllTable(f func(ch chan<- dao.BiliVideo)) {
	ch := make(chan dao.BiliVideo, 100)
	go f(ch)
	idMap := make(map[string]struct{}, 40000)
	for {
		video, ok := <-ch
		if !ok {
			break
		}
		idMap[video.Id] = struct{}{}
	}
	fmt.Println(len(idMap))
}

func TestReadAllTable1(t *testing.T) {
	testReadAllTable(dao.ReadAllTable1) //ReadAllTable1 use time 173 ms
}

// 2024/10/13 18:15:20 S:/VsCode_Repo/GoLandCode/RADIC/course/dao/bili_video.go:203
// [18.846ms] [rows:2710] SELECT * FROM `bili_video`
// 2024/10/13 18:15:20 ReadAllTable1 read 2710 records
// ReadAllTable1 use time 40 ms
// 2710

func TestReadAllTable2(t *testing.T) {
	testReadAllTable(dao.ReadAllTable2) //ReadAllTable2 use time 2654 ms
}

// 2024/10/13 18:16:04 S:/VsCode_Repo/GoLandCode/RADIC/course/dao/bili_video.go:227
// [2.476ms] [rows:0] SELECT * FROM `bili_video` LIMIT 500 OFFSET 2710
// 2024/10/13 18:16:04 ReadAllTable2 read 2710 records
// ReadAllTable2 use time 58 ms
// 2710

func TestReadAllTable3(t *testing.T) {
	testReadAllTable(dao.ReadAllTable3) //ReadAllTable3 use time 262 ms
}

// 2024/10/13 18:16:55 S:/VsCode_Repo/GoLandCode/RADIC/course/dao/bili_video.go:259
// [0.498ms] [rows:0] SELECT * FROM `bili_video` WHERE id>'BV1zz4y1x7Zm' LIMIT 500
// 2024/10/13 18:16:55 ReadAllTable2 read 2710 records
// ReadAllTable3 use time 49 ms
// 2710

// go test -v ./course/dao/test -run=^TestDumpDataFromFile2DB1$ -count=1
// go test -v ./course/dao/test -run=^TestDumpDataFromFile2DB2$ -count=1
// go test -v ./course/dao/test -run=^TestDumpDataFromFile2DB3$ -count=1
// go test -v ./course/dao/test -run=^TestReadAllTable1$ -count=1 -timeout=30m
// go test -v ./course/dao/test -run=^TestReadAllTable2$ -count=1 -timeout=30m
// go test -v ./course/dao/test -run=^TestReadAllTable3$ -count=1 -timeout=30m
