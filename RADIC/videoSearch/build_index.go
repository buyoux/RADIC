package videoSearch

import (
	indexservice "RADIC/index_service"
	"RADIC/types"
	"RADIC/util"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	farmhash "github.com/leemcloughlin/gofarmhash"
)

func BuildIndexFromFile(csvFile string, indexer *indexservice.Indexer, totalWorkers, workerIndex int) {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Printf("open file %s failed: %s", csvFile, err)
	}
	defer file.Close()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	reader := csv.NewReader(file) // 读取csv文件
	progress := 0
	for {
		record, err := reader.Read() // 读取csv文件的一行，record是个切片
		if err != nil {
			if err != io.EOF {
				log.Printf("read record failed: %s", err)
			}
			break
		}

		if len(record) < 10 {
			continue
		}
		docId := strings.TrimPrefix(record[0], "https://www.bilibili.com/video/")
		//只用一部分的视频数据
		if totalWorkers > 0 && int(farmhash.Hash32WithSeed([]byte(docId), 0))%totalWorkers != workerIndex {
			continue
		}
		video := &BiliVideo{
			Id:     docId,
			Title:  record[1],
			Author: record[3],
		}
		if len(record[2]) > 4 {
			t, err := time.ParseInLocation("2006/1/2 15:4", record[2], loc)
			if err != nil {
				log.Printf("parse time %s failed: %s", record[2], err)
			} else {
				video.PostTime = t.Unix()
			}
		}
		n, _ := strconv.Atoi(record[4])
		video.View = int32(n)
		n, _ = strconv.Atoi(record[5])
		video.Like = int32(n)
		n, _ = strconv.Atoi(record[6])
		video.Coin = int32(n)
		n, _ = strconv.Atoi(record[7])
		video.Favorite = int32(n)
		n, _ = strconv.Atoi(record[8])
		video.Share = int32(n)
		keywords := strings.Split(record[9], ",")
		if len(keywords) > 0 {
			for _, word := range keywords {
				word = strings.TrimSpace(word)
				if len(word) > 0 {
					video.Keywords = append(video.Keywords, strings.ToLower(word))
				}
			}
		}
		// 构建BiliVideo实体后，写入索引
		AddVideo2Index(video, indexer)
		progress++
		if progress%100 == 0 {
			log.Printf("read record progress = %d\n", progress)
		}
	}
	util.Log.Printf("add %d documents to index totally", progress)
}

// 把一条视频信息写入索引(可能是create，也可能是update)
// 实时更新索引时可调该函数
func AddVideo2Index(video *BiliVideo, indexer *indexservice.Indexer) {
	doc := types.Document{Id: video.Id}
	bs, err := proto.Marshal(video)
	if err == nil {
		doc.Bytes = bs
	} else {
		log.Printf("serielize video failed: %s", err)
		return
	}
	keywords := make([]*types.Keyword, 0, len(video.Keywords))
	for _, word := range video.Keywords {
		keywords = append(keywords, &types.Keyword{Field: "content", Word: strings.ToLower(word)})
	}
	if len(video.Author) > 0 {
		keywords = append(keywords, &types.Keyword{Field: "author", Word: strings.ToLower(strings.TrimSpace(video.Author))})
	}
	doc.Keywords = keywords
	doc.BitsFeature = GetClassBits(video.Keywords)

	indexer.AddDoc(doc)
}
