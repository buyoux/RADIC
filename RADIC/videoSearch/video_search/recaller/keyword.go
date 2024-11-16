package recaller

import (
	"RADIC/types"
	"RADIC/videoSearch"
	"RADIC/videoSearch/video_search/common"
	"strings"

	"github.com/gogo/protobuf/proto"
)

type KeywordRecaller struct{}

func (KeywordRecaller) Recall(ctx *common.VideoSearchContext) []*videoSearch.BiliVideo {
	request := ctx.Request
	if request == nil {
		return nil
	}
	indexer := ctx.Indexer
	if indexer == nil {
		return nil
	}

	keywords := request.Keywords
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
	docs := indexer.Search(query, 0, 0, orFlags)
	videos := make([]*videoSearch.BiliVideo, 0, len(docs))
	for _, doc := range docs {
		var video videoSearch.BiliVideo
		if err := proto.Unmarshal(doc.Bytes, &video); err == nil {
			videos = append(videos, &video)
		}
	}
	return videos
}
