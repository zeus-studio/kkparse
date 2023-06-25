package parser

import (
	"errors"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
)

type weiShi struct {
}

func (w weiShi) parseVideoID(videoId string) (*VideoParseInfo, error) {
	reqUrl := "https://h5.weishi.qq.com/webapp/json/weishi/WSH5GetPlayPage?feedid=" + videoId
	client := resty.New()
	if ProxyIP != "" {
		client.SetProxy("http://" + ProxyIP)
	}
	res, err := client.R().
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(reqUrl)
	if err != nil {
		return nil, err
	}

	data := gjson.GetBytes(res.Body(), "data.feeds.0")
	author := data.Get("poster.nick").String()
	avatar := data.Get("poster.avatar").String()
	title := data.Get("feed_desc_withat").String()
	videoUrl := data.Get("video_url").String()
	cover := data.Get("images.0.url").String()

	parseRes := &VideoParseInfo{
		Title:    title,
		VideoUrl: videoUrl,
		CoverUrl: cover,
		Source:   SourceWeiShi,
	}
	parseRes.Author.Name = author
	parseRes.Author.Avatar = avatar

	return parseRes, nil
}

func (w weiShi) parseShareUrl(shareUrl string) (*VideoParseInfo, error) {
	urlRes, err := url.Parse(shareUrl)
	if err != nil {
		return nil, err
	}

	videoId := urlRes.Query().Get("id")
	if len(videoId) <= 0 {
		return nil, errors.New("parse video_id from share url fail")
	}

	return w.parseVideoID(videoId)
}
