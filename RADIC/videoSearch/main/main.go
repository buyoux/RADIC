package main

import (
	"RADIC/internal/kvdb"
	"RADIC/util"
	"RADIC/videoSearch/handler"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	mode         = flag.Int("mode", 1, "启动哪类服务。1-standalong web server, 2-grpc index server, 3-distrubuted index server")
	rebuildIndex = flag.Bool("index", false, "server启动时是否需要重建索引")
	port         = flag.Int("port", 0, "server的工作端口")
	dbPath       = flag.String("dbPath", "", "正排索引数据的存放路径")
	totalWorkers = flag.Int("totalWorkers", 0, "分布式环境中一共有几台index worker")
	workerIndex  = flag.Int("workerIndex", 0, "本机是第几台index worker(从0开始编号)")
)

var (
	dbType      = kvdb.BOLT                                    // 正排索引使用哪种KV数据库
	csvFile     = util.ProjectRootPath + "data/bili_video.csv" // 原始的数据文件，由它来创建索引
	etcdServers = []string{"127.0.0.1:2379"}                   // etcd集群的地址
)

func StartGin() {
	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	engine.Static("js", "videoSearch/views/js")
	engine.Static("css", "videoSearch/views/css")
	engine.Static("img", "videoSearch/views/img")
	//在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	engine.StaticFile("/favicon.ico", "img/fav.jpg")
	//使用这些.html文件时就不需要加路径了
	engine.LoadHTMLFiles("videoSearch/views/search.html", "videoSearch/views/up_search.html")

	engine.Use(handler.GetUserInfo)
	classes := [...]string{"资讯", "社会", "热点", "生活", "知识", "环球", "游戏", "综合", "日常", "影视", "科技", "编程"} //数组，非切片
	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "search.html", classes)
	})
	engine.GET("/up", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "up_search.html", classes)
	})

	engine.POST("/search", handler.SearchAll)
	engine.POST("/up_search", handler.SearchByAuthor)

	engine.Run("127.0.0.1:" + strconv.Itoa(*port))
}

func main() {
	flag.Parse()
	fmt.Println(*totalWorkers)
	switch *mode {
	case 1, 3:
		// 1：单机模式，索引功能嵌套在web server内部。
		// 3：分布式模式，web server内持有一个哨兵，通过哨兵去访问各个grpc index server
		WebServerMain(*mode)
		StartGin()
	case 2:
		GrpcIndexerMain()
	}
}

// 单机模式
// go run ./videoSearch/main -mode=1 -index=true -port=5678 -dbPath=data/local_db/video_bolt
// go run ./videoSearch/main -mode=1 -index=false -port=5678 -dbPath=data/local_db/video_bolt

// 启动Grpc Server分布式模式
// go run ./videoSearch/main -mode=2 -index=true -port=5600 -dbPath=data/local_db/video_bolt -totalWorkers=2 -workerIndex=0
// go run ./videoSearch/main -mode=2 -index=true -port=5601 -dbPath=data/local_db/video_bolt -totalWorkers=2 -workerIndex=1

// 启动Web Server
// go run ./videoSearch/main -mode=3 -index=false -port=5678
