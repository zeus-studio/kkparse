package parser

import (
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type xHongShu struct {
}

func (x xHongShu) parseShareUrl(shareUrl string) (*VideoParseInfo, error) {

	// 前置请求，获取重定向链接
	prevClient := resty.New()
	if ProxyIP != "" {
		prevClient.SetProxy("http://" + ProxyIP)
	}
	prevClient.SetRedirectPolicy(resty.NoRedirectPolicy())
	_, resErr := prevClient.R().Get(shareUrl)
	if resErr == nil {
		return nil, fmt.Errorf("parse video id from share url fail")
	}
	errStr := resErr.Error()
	urlMatches := regexp.MustCompile(`"([\w\W]*)"`).FindAllStringSubmatch(errStr, -1)
	truthUrl := urlMatches[0][1]

	client := resty.New()
	if ProxyIP != "" {
		client.SetProxy("http://" + ProxyIP)
	}
	res, err := client.R().
		SetHeader(HttpHeaderUserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36").
		SetHeader(HttpHeaderCookie, `webBuild=2.11.5; xsecappid=xhs-pc-web; a1=188e69400cau1r372yxyqk2y4ecch1m485g973ma530000298254; webId=ca2347778834bf0d9bda1237ec99405e; gid=yYYdKj48yiYKyYYdKj4882q6S07yFqWJvCvAk02TJv4dSjq8Sxy64W888JjYJ248ddyqDSYq; gid.sign=RpF8W4B7QbYFiIDHn8sFQOdbI0s; web_session=030037a3a677ac37ce329465ed234a4dca3449; websectiga=634d3ad75ffb42a2ade2c5e1705a73c845837578aeb31ba0e442d75c648da36a; sec_poison_id=0899751f-9d40-4d4a-b297-b1cbad491e55`).
		SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		SetHeader("accept-language", "zh-CN,zh;q=0.9").
		SetHeader("cache-control", "no-cache").
		SetHeader("pragma", "no-cache").
		SetHeader("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`).
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", `"macOS"`).
		SetHeader("sec-fetch-dest", "document").
		SetHeader("sec-fetch-mode", "navigate").
		SetHeader("sec-fetch-site", "none").
		SetHeader("sec-fetch-user", "?1").
		SetHeader("upgrade-insecure-requests", "1").
		SetHeader("authority", "www.xiaohongshu.com").
		Get(truthUrl)

	if err != nil {
		return nil, err
	}

	body := string(res.Body())
	pattern := `(window\.__INITIAL_STATE__\s?=\s?)(\{[\w\W]*)(\}</script>)`
	re := regexp.MustCompile(pattern)

	// 提取匹配结果
	matches := re.FindAllStringSubmatch(body, -1)

	if matches == nil || len(matches[0]) == 0 {
		return nil, fmt.Errorf("解析失败")
	}

	// 获取 backupUrls 的 Json 字符串
	data := gjson.Parse(matches[0][2] + "}")

	videoInfo := &VideoParseInfo{
		Title:    data.Get("note.note.title").String(),
		MusicUrl: "",
		CoverUrl: "",
		Source:   SourceXiaoHongShu,
	}
	videoInfo.Author.Uid = data.Get("note.note.user.userId").String()
	videoInfo.Author.Name = data.Get("note.note.user.nickname").String()
	videoInfo.Author.Avatar = data.Get("note.note.user.avatar").String()

	var videoUrl = data.Get("note.note.video.media.stream.h265.0.backupUrls.0").String()
	if videoUrl == "" {
		videoUrl = data.Get("note.note.video.media.stream.av1.0.backupUrls.0").String()
	}
	if videoUrl == "" {
		videoUrl = data.Get("note.note.video.media.stream.h264.0.backupUrls.0").String()
	}
	videoInfo.VideoUrl = videoUrl

	return videoInfo, nil
}
