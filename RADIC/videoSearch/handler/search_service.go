package handler

import (
	indexservice "RADIC/index_service"
	"RADIC/types"
	"RADIC/util"
	"RADIC/videoSearch"
	videosearch "RADIC/videoSearch/video_search"
	"RADIC/videoSearch/video_search/common"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
)

var Indexer indexservice.IIndexer

func clearnKeywords(words []string) []string {
	keywords := make([]string, 0, len(words))
	for _, w := range words {
		word := strings.TrimSpace(strings.ToLower(w))
		if len(word) > 0 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}

func Search(ctx *gin.Context) {
	var request videoSearch.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	keywords := clearnKeywords(request.Keywords)
	if len(keywords) == 0 && len(request.Author) == 0 {
		ctx.String(http.StatusBadRequest, "关键词和作者不能同时为空")
		return
	}

	query := new(types.TermQuery)
	if len(keywords) > 0 {
		for _, word := range keywords {
			query = query.And(types.NewTermQuery("content", word))
		}
	}

	if len(request.Author) > 0 {
		query = query.And(types.NewTermQuery("author", strings.ToLower(request.Author)))
	}

	orFlags := []uint64{videoSearch.GetClassBits(request.Classes)}
	docs := Indexer.Search(query, 0, 0, orFlags)
	videos := make([]videoSearch.BiliVideo, 0, len(docs))
	for _, doc := range docs {
		var video videoSearch.BiliVideo
		if err := proto.Unmarshal(doc.Bytes, &video); err == nil {
			if video.View >= int32(request.ViewFrom) && (request.ViewTo <= 0 || video.View <= int32(request.ViewTo)) {
				videos = append(videos, video)
			}
		}
	}
	util.Log.Printf("return %d videos", len(videos))
	ctx.JSON(http.StatusOK, videos)
}

func SearchAll(ctx *gin.Context) {
	var request videoSearch.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	keywords := clearnKeywords(request.Keywords)
	if len(keywords) == 0 && len(request.Author) == 0 {
		ctx.String(http.StatusBadRequest, "关键词和作者不能同时为空")
		return
	}

	searchCtx := &common.VideoSearchContext{
		Ctx:     context.Background(),
		Request: &request,
		Indexer: Indexer,
	}
	searcher := videosearch.NewAllVideoSearch()
	videos := searcher.Search(searchCtx)

	util.Log.Printf("return %d videos", len(videos))
	ctx.JSON(http.StatusOK, videos)
}

func SearchByAuthor(ctx *gin.Context) {
	var request videoSearch.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	keywords := clearnKeywords(request.Keywords)
	if len(keywords) == 0 && len(request.Author) == 0 {
		ctx.String(http.StatusBadRequest, "关键词和作者不能同时为空")
		return
	}

	userName, ok := ctx.Value("user_name").(string) // 从gin.Context里取得userName
	if !ok || len(userName) == 0 {
		ctx.String(http.StatusBadRequest, "获取不到登录用户名")
		return
	}
	searchCtx := &common.VideoSearchContext{
		Ctx:     context.WithValue(context.Background(), common.UN("user_name"), userName),
		Request: &request,
		Indexer: Indexer,
	}
	searcher := videosearch.NewUpVideoSearch()
	videos := searcher.Search(searchCtx)

	util.Log.Printf("return %d videos", len(videos))
	ctx.JSON(http.StatusOK, videos)
}
