package parser

import (
	"errors"

	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
)

type huoShan struct {
}

func (h huoShan) parseVideoID(videoId string) (*VideoParseInfo, error) {
	reqUrl := "https://share.huoshan.com/api/item/info?item_id=" + videoId
	client := resty.New()
	res, err := client.R().
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(reqUrl)
	if err != nil {
		return nil, err
	}

	data := gjson.GetBytes(res.Body(), "data.item_info")
	videoUrl := data.Get("url").String()
	cover := data.Get("cover").String()

	parseRes := &VideoParseInfo{
		VideoUrl: videoUrl,
		CoverUrl: cover,
		Source:   SourceHuoShan,
	}

	return parseRes, nil
}

func (h huoShan) parseShareUrl(shareUrl string) (*VideoParseInfo, error) {

	client := resty.New()
	if ProxyIP != "" {
		client.SetProxy("http://" + ProxyIP)
	}
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	res, _ := client.R().
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(shareUrl)
	// 这里会返回err, auto redirect is disabled

	locationRes, err := res.RawResponse.Location()
	if err != nil {
		return nil, err
	}

	videoId := locationRes.Query().Get("item_id")
	if len(videoId) <= 0 {
		return nil, errors.New("parse video id from share url fail")
	}

	return h.parseVideoID(videoId)
}
