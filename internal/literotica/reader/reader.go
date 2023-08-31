// Package reader 读
package reader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/r0oot/tools/internal/literotica/processor"
	"github.com/r0oot/tools/internal/literotica/protocol"

	htmlstrip "github.com/grokify/html-strip-tags-go"
)

type Reader struct {
	basicInfo *protocol.BasicInfo
	cursor    int
}

func (r *Reader) GetInfoAndGetter(url string) (*protocol.BasicInfo, processor.NextSectionGetter, error) {
	basicInfo, err := getBasicInfo(url)
	if err != nil {
		return nil, nil, err
	}
	r.basicInfo = basicInfo
	return basicInfo, r.GetSectionInfo, nil
}

// GetSectionInfo 取章节内容
func (r *Reader) GetSectionInfo() (*protocol.SectionInfo, error) {
	if r.cursor >= len(r.basicInfo.IndexURL) {
		return nil, nil
	}
	url := r.basicInfo.IndexURL[r.cursor]
	s, err := r.getSectionPageByPage(url, true)
	if err != nil {
		return nil, err
	}
	r.cursor++
	return s, nil
}

func (r *Reader) getSectionPageByPage(url string, subpage bool) (*protocol.SectionInfo, error) {
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
	{
		reg := regexp.MustCompile(`<a class="l_bJ l_bL" title="Next Page" href="([^"]*?)"`)
		result := reg.FindAllStringSubmatch(body, -1)
		for _, pageR := range result {
			if len(pageR) < 2 {
				continue
			}
			pageCnt, err := r.getSectionPageByPage("https://www.literotica.com"+pageR[1], false)
			if err != nil {
				return nil, errors.New("page 获取失败:" + pageR[1])
			}
			cnt += pageCnt.Content
		}
	}

	return &protocol.SectionInfo{
		Title:   title,
		Content: htmlstrip.StripTags(cnt),
	}, nil
}

// getBasicInfo 获取索引页信息 返回为title,[]indexurl,error
func getBasicInfo(url string) (*protocol.BasicInfo, error) {
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
	return &protocol.BasicInfo{
		Title:    title,
		IndexURL: index,
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
