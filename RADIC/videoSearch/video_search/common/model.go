package common

import (
	indexservice "RADIC/index_service"
	"RADIC/videoSearch"
	"context"
)

type VideoSearchContext struct {
	Ctx     context.Context            // 上下文参数
	Indexer indexservice.IIndexer      // 索引。既可能是本地的Indexer,也可能是分布式的哨兵Sentinel
	Request *videoSearch.SearchRequest // 搜索请求
	Videos  []*videoSearch.BiliVideo   // 搜索结果
}

type UN string
