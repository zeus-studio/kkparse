package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// ParseVideoShareUrl 根据视频分享链接解析视频信息
func ParseVideoShareUrl(shareUrl string) (*VideoParseInfo, error) {
	// 根据分享url判断source
	source := ""
	for itemSource, itemSourceInfo := range videoSourceInfoMapping {
		for _, itemUrlDomain := range itemSourceInfo.VideoShareUrlDomain {
			if strings.Contains(shareUrl, itemUrlDomain) {
				source = itemSource
				break
			}
		}
		if len(source) > 0 {
			break
		}
	}

	// 没有找到对应source
	if len(source) <= 0 {
		return nil, fmt.Errorf("share url [%s] not have source config", shareUrl)
	}

	// 没有对应的视频链接解析方法
	urlParser := videoSourceInfoMapping[source].VideoShareUrlParser
	if urlParser == nil {
		return nil, fmt.Errorf("source %s has no video share url parser", source)
	}

	// 获取 IP 地址
	getProxyIp()

	return urlParser.parseShareUrl(shareUrl)
}

// ParseVideoId 根据视频id解析视频信息
func ParseVideoId(source, videoId string) (*VideoParseInfo, error) {
	if len(videoId) <= 0 || len(source) <= 0 {
		return nil, errors.New("video id or source is empty")
	}

	idParser := videoSourceInfoMapping[source].VideoIdParser
	if idParser == nil {
		return nil, fmt.Errorf("source %s has no video id parser", source)
	}

	return idParser.parseVideoID(videoId)
}

// BatchParseVideoId 根据视频id批量解析视频信息
func BatchParseVideoId(source string, videoIds []string) (map[string]BatchParseItem, error) {
	if len(videoIds) <= 0 || len(source) <= 0 {
		return nil, errors.New("videos id or source is empty")
	}

	idParser := videoSourceInfoMapping[source].VideoIdParser
	if idParser == nil {
		return nil, fmt.Errorf("source %s has no video id parser", source)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	parseMap := make(map[string]BatchParseItem, len(videoIds))
	for _, v := range videoIds {
		wg.Add(1)
		videoId := v
		go func(videoId string) {
			defer wg.Done()

			parseInfo, parseErr := ParseVideoId(source, videoId)
			mu.Lock()
			parseMap[videoId] = BatchParseItem{
				ParseInfo: parseInfo,
				Error:     parseErr,
			}
			mu.Unlock()
		}(videoId)
	}
	wg.Wait()

	return parseMap, nil
}

var ProxyPoolList []string = make([]string, 30)
var ProxyIP string

func getProxyIp() {
	client := resty.New()
	resp, err := client.R().
		Get("https://api.docip.net/v1/get_openproxy?api_key=4TBDalUzvzTlXePFwNWOY64843c1b&num=30&proxy_type=1&country_type=1&sort_type=2&ports=&unports=&areas=&unareas=&quchong=0&format=json")
	if err != nil {
		ProxyIP = ProxyPoolList[0]
		return
	}

	if !json.Valid(resp.Body()) {
		ProxyIP = ""
		return
	}
	data := gjson.ParseBytes(resp.Body()).Array()
	ProxyIP = data[0].String()
	ProxyPoolList = ProxyPoolList[:0]
	for _, v := range data {
		ProxyPoolList = append(ProxyPoolList, v.String())
	}
}
