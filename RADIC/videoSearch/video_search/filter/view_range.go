package filter

import (
	"RADIC/videoSearch"
	"RADIC/videoSearch/video_search/common"
)

type ViewFilter struct{}

func (ViewFilter) Apply(ctx *common.VideoSearchContext) {
	request := ctx.Request
	if request == nil {
		return
	}
	if request.ViewFrom >= request.ViewTo {
		return
	}
	videos := make([]*videoSearch.BiliVideo, 0, len(ctx.Videos))
	for _, video := range ctx.Videos {
		if video.View >= int32(request.ViewFrom) && video.View <= int32(request.ViewTo) {
			videos = append(videos, video)
		}
	}
	ctx.Videos = videos
}
