package dao

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var loc *time.Location

const BATCH_SIZE = 300

// 适合使用init()的典型场景：全局变量的初始化放到init()里，且没有任何前提依赖
// 如数据库连接可能会有其他日志依赖，可以放到Init()中设计依赖顺序
func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}

type BiliVideo struct {
	Id       string //结构体里的驼峰转为蛇形，即mysql表里的列名
	Title    string
	Author   string
	PostTime time.Time
	Keywords string
	View     int
	ThumbsUp int
	Coin     int
	Favorite int
	Share    int
}

func (BiliVideo) TableName() string {
	return "bili_video" //指定表名
}

func parseFileLine(record []string) *BiliVideo {
	video := &BiliVideo{
		Title:  record[1],
		Author: record[3],
	}
	urlPaths := strings.Split(record[0], "/")
	video.Id = urlPaths[len(urlPaths)-1]
	if len(record[2]) > 4 {
		t, err := time.ParseInLocation("2006/1/2 15:4", record[2], loc)
		if err != nil {
			log.Printf("parse time %s failed: %s", record[2], err)
		} else {
			video.PostTime = t
		}
	}
	n, _ := strconv.Atoi(record[4])
	video.View = n
	n, _ = strconv.Atoi(record[5])
	video.ThumbsUp = n
	n, _ = strconv.Atoi(record[6])
	video.Coin = n
	n, _ = strconv.Atoi(record[7])
	video.Favorite = n
	n, _ = strconv.Atoi(record[8])
	video.Share = n
	video.Keywords = strings.ToLower(record[9])
	return video
}

func readFile(csvFile string, ch chan<- *BiliVideo) {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Printf("open file %s failed: %s", csvFile, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file) // 读取CSV文件
	for {
		// 读取csv文件的一行，record是一个切片
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				log.Printf("read record failed: %s", err)
			}
			break
		}
		// 避免数组越界，发生panic
		if len(record) < 10 {
			continue
		}
		video := parseFileLine(record)
		ch <- video
	}
	close(ch) // 生产方读取数据结束后，一定要关闭管道close channel
}

func DumpDataFromFile2DB1(csvFile string) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("DumpDataFromFile2DB1 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)
	ch := make(chan *BiliVideo, 200)
	go readFile(csvFile, ch)

	db := GetSearchDBConnection()
	for {
		video, ok := <-ch
		if !ok {
			break
		}
		err := db.Create(video).Error
		checkErr(err)
	}
}

func DumpDataFromFile2DB2(csvFile string) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("DumpDataFromFile2DB1 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)
	ch := make(chan *BiliVideo, 200)
	go readFile(csvFile, ch)

	db := GetSearchDBConnection()
	tx := db.Begin()
	i := 0
	for {
		video, ok := <-ch
		if !ok {
			break
		}
		// 通过事务提交insert请求
		tx.Create(video)
		i++
		if i >= BATCH_SIZE {
			// 300次insert提交一次事务
			err := tx.Commit().Error
			checkErr(err)
			// 不能在一个事务上重复commit，需要新开一个
			tx = db.Begin()
			i = 0
		}
	}
	err := tx.Commit().Error
	checkErr(err)
}

func DumpDataFromFile2DB3(csvFile string) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("DumpDataFromFile2DB1 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)
	ch := make(chan *BiliVideo, 200)
	go readFile(csvFile, ch)

	db := GetSearchDBConnection()
	buffer := make([]*BiliVideo, 0, BATCH_SIZE)
	for {
		video, ok := <-ch
		if !ok {
			break
		}
		buffer = append(buffer, video)
		if len(buffer) >= BATCH_SIZE {
			err := db.CreateInBatches(buffer, BATCH_SIZE).Error
			checkErr(err)
			buffer = make([]*BiliVideo, 0, BATCH_SIZE)
		}
	}
	err := db.CreateInBatches(buffer, BATCH_SIZE).Error
	checkErr(err)
}

func checkErr(err error) {
	// et := reflect.TypeOf(err).Elem()
	// fmt.Println(et, et.PkgPath(), et.Name())
	var sqlErr *mysql.MySQLError
	if errors.As(err, &sqlErr) {
		if sqlErr.Number != 1062 {
			panic(err)
		}
	}
}

func ReadAllTable1(ch chan<- BiliVideo) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("ReadAllTable1 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)

	db := GetSearchDBConnection()
	var datas []BiliVideo
	// select * from bili_video;  禁止这种写法，绝对是慢查询
	// 通常超过100ms定义为慢查询
	if err := db.Select("*").Find(&datas).Error; err != nil {
		log.Printf("ReadAllTable1 failed: %s", err)
	}
	for _, data := range datas {
		ch <- data
	}
	log.Printf("ReadAllTable1 read %d records", len(datas))
	close(ch)
}

func ReadAllTable2(ch chan<- BiliVideo) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("ReadAllTable2 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)

	db := GetSearchDBConnection()
	offset := 0
	const BATCH = 500
	for {
		t0 := time.Now()
		var datas []BiliVideo
		// select * from bili_video limit offset, BATCH;
		// 实际上执行limit 0, offset + BATCH
		if err := db.Select("*").Offset(offset).Limit(BATCH).Find(&datas).Error; err != nil {
			log.Printf("ReadAllTable2 failed: %s", err)
			break
		} else {
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				ch <- data
			}
			offset += len(datas)
		}
		fmt.Printf("offset=%d use time %dms\n", offset, time.Since(t0).Milliseconds())
	}
	log.Printf("ReadAllTable2 read %d records", offset)
	close(ch)
}

func ReadAllTable3(ch chan<- BiliVideo) {
	begin := time.Now()
	defer func(begin time.Time) {
		fmt.Printf("ReadAllTable3 use time %d ms\n", time.Since(begin).Milliseconds())
	}(begin)

	db := GetSearchDBConnection()
	maxid := ""
	const BATCH = 500
	total := 0
	for {
		t0 := time.Now()
		var datas []BiliVideo
		// select * from bili_video where id>maxid limit BATCH;
		if err := db.Select("*").Where("id>?", maxid).Limit(BATCH).Find(&datas).Error; err != nil {
			log.Printf("ReadAllTable3 failed: %s", err)
			break
		} else {
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				ch <- data
			}
			maxid = datas[len(datas)-1].Id
			total += len(datas)
		}
		fmt.Printf("total=%d use time %dms\n", total, time.Since(t0).Milliseconds())
	}
	log.Printf("ReadAllTable2 read %d records", total)
	close(ch)
}
