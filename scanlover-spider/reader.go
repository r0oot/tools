package main

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"time"
)

type content struct {
	images   []string
	nextPage string
}

var (
	imgReg      = regexp.MustCompile(`<img src="([^"]*?)"`)
	nextPageReg = regexp.MustCompile(`<a href="([^"]*?)">Next Page`)
)

// getAndParse 取页面中的全部图片
func getAndParse(url string) (*content, error) {
	// 下载一个页面
	body, err := download(url)
	if err != nil {
		return nil, err
	}
	ret := &content{}
	// 图片链接
	result := imgReg.FindAllStringSubmatch(body, -1)
	for _, one := range result {
		if len(one) < 2 {
			return nil, errors.New("解析图片错误")
		}
		ret.images = append(ret.images, one[1])
	}

	// 下一页
	r := nextPageReg.FindStringSubmatch(body)
	if len(r) >= 2 {
		ret.nextPage = html.UnescapeString(r[1])
	}
	return ret, nil
}

// download 下载一个页面
func download(url string) (string, error) {
	time.Sleep(1 * time.Second)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("[url:%+v\t\t\t size:%+v]\n", url, len(body))
	return string(body), nil
}
