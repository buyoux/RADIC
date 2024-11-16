package main

import (
	indexservice "RADIC/index_service"
	"RADIC/videoSearch"
	"RADIC/videoSearch/handler"
	"os"
	"os/signal"
	"syscall"
)

func WebServerInit(mode int) {
	switch mode {
	case 1:
		standalongIndexer := new(indexservice.Indexer)
		if err := standalongIndexer.Init(5000, dbType, *dbPath); err != nil {
			panic(err)
		}
		if *rebuildIndex {
			videoSearch.BuildIndexFromFile(csvFile, standalongIndexer, 0, 0) // 重建索引
		} else {
			standalongIndexer.LoadFromIndexFile() // 直接从持久化的正排索引文件里加载
		}
		handler.Indexer = standalongIndexer
	case 3:
		handler.Indexer = indexservice.NewSentinel(etcdServers)
	default:
		panic("invalid mode")
	}
}

func WebServerTeardown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	handler.Indexer.Close() // 接收到kill信号时关闭索引
	os.Exit(0)              // 然后退出程序
}

func WebServerMain(mode int) {
	go WebServerTeardown()
	WebServerInit(mode)
}
