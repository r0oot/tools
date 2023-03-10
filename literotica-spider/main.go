package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	htmlstrip "github.com/grokify/html-strip-tags-go"
)

var (
	url      string
	offset   int
	bookType string
)

func init() {
	flag.StringVar(&url, "url", "", "填文章的链接")
	flag.StringVar(&bookType, "booktype", "epub", "支持txt和epub")
	flag.IntVar(&offset, "offset", 0, "要跳过几个文章") // TODO
	flag.Parse()
}

func main() {
	// 获取索引
	info, err := getIndex(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	var writer writer
	if bookType == "txt" {
		writer = &txtWriter{}
	}
	if bookType == "epub" {
		writer = &epubWriter{}
	}
	// 创建文件句柄
	writer.Create(info.title)
	defer writer.Save()

	// 遍历读取并写入
	for _, indexURL := range info.indexURL {
		indexInfo, err := getIndexContent(indexURL, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		indexInfo.title = strings.ReplaceAll(indexInfo.title, info.title, "")
		writer.AddSection(indexInfo.title, indexInfo.content)
	}
}

type info struct {
	title    string
	indexURL []string
}

type indexInfo struct {
	title   string
	content string
}

// getIndexContent 获取一个索引的内容
func getIndexContent(url string, subpage bool) (*indexInfo, error) {
	body, err := download(url)
	if err != nil {
		return nil, err
	}
	// 标题
	var title string
	if subpage {
		reg := regexp.MustCompile(`<h1 class="j_bm headline j_eQ">(.*?)<\/h1>`)
		result := reg.FindStringSubmatch(body)
		if len(result) < 2 {
			return nil, errors.New("获取索引标题失败，请检查页面内容和正则")
		}
		title = result[1]
	}

	// 内容
	reg := regexp.MustCompile(`<div class="panel article aa_eQ"><div class="aa_ht"><div>(.*?)</div></div><div class="aa_ht"></div>`)
	result := reg.FindStringSubmatch(body)
	if len(result) < 2 {
		return nil, errors.New("获取索引内容失败，请检查页面内容和正则")
	}
	cnt := result[1]
	cnt = strings.ReplaceAll(cnt, "<p>", "")
	cnt = strings.ReplaceAll(cnt, "</p>", "\n\n")

	// 分页
	if subpage {
		reg := regexp.MustCompile(`<a class="l_bJ" href="([^"]*?)"`)
		result := reg.FindAllStringSubmatch(body, -1)
		for _, pageR := range result {
			if len(pageR) < 2 {
				continue
			}
			pageCnt, err := getIndexContent("https://www.literotica.com"+pageR[1], false)
			if err != nil {
				return nil, errors.New("page 获取失败:" + pageR[1])
			}
			cnt += pageCnt.content
		}
	}
	return &indexInfo{
		title:   title,
		content: htmlstrip.StripTags(cnt),
	}, nil
}

// getIndex 获取信息 返回为title,[]indexurl,error
func getIndex(url string) (*info, error) {
	// 下载一个页面
	body, err := download(url)
	if err != nil {
		return nil, err
	}
	// 在页面里找这个系列的链接
	reg := regexp.MustCompile(`href="(([^"]*?)\/series\/([^"]*)?)"`)
	result := reg.FindStringSubmatch(body)
	if len(result) < 2 {
		return nil, errors.New("获取系列的链接失败，请检查页面内容和正则")
	}
	seriesURL := result[1]
	fmt.Printf("[==获取series链接成功:%+v]\n", seriesURL)
	// 打开系列的链接，获取title和全部索引
	seriesBody, err := download(seriesURL)
	if err != nil {
		return nil, err
	}
	// 匹配标题
	titleReg := regexp.MustCompile(`<h1 class="j_bm headline">(.*?)</h1>`)
	result = titleReg.FindStringSubmatch(seriesBody)
	if len(result) < 2 {
		return nil, errors.New("获取系列的标题失败，请检查页面内容和正则")
	}
	title := result[1]
	// 匹配索引的链接
	indexReg := regexp.MustCompile(`<a href="([^"]*?)" class="br_rj">`)
	indexResult := indexReg.FindAllStringSubmatch(seriesBody, -1)
	var index []string
	for _, r := range indexResult {
		if len(r) < 2 {
			return nil, errors.New("获取索引的链接失败，请检查页面内容和正则")
		}
		index = append(index, r[1])
	}
	return &info{
		title:    title,
		indexURL: index,
	}, nil
}

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
