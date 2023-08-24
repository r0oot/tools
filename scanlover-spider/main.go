package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	url  string
	dir  string
	page int
)

func init() {
	flag.StringVar(&url, "url", "https://scanlover.com/d/1982?id=1982-momo-sakura&near=&page=41", "填文章的链接")
	flag.StringVar(&dir, "dir", "", "")
	flag.IntVar(&page, "page", 1, "")
	flag.Parse()
}

func main() {
	for {
		cnt, err := getAndParse(url)
		if err != nil {
			fmt.Println("下载解析失败", err)
			return
		}
		index := 1
		for _, imgURL := range cnt.images {
			// gif 格式的不要
			ext := strings.ToLower(filepath.Ext(imgURL))
			if ext == ".gif" {
				continue
			}

			fmt.Println(page, index, imgURL)
			content, err := download(imgURL)
			if err != nil {
				fmt.Println("图片下载失败", err)
				continue
			}
			fileName := fmt.Sprintf("%s/%d_%d%s", dir, page, index, ext)
			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println("文件打开失败", err)
				return
			}
			writer := bufio.NewWriter(file)
			writer.WriteString(content)
			writer.Flush()
			file.Close()
			index++
		}
		page++
		url = cnt.nextPage
		if url == "" {
			fmt.Println("finished")
			return
		}
	}
}
