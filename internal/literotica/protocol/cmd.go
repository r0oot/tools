package protocol

import "flag"

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
