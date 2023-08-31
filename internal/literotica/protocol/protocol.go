// Package 协议
package protocol

// BookType 电子书类型
type BookType string

const (
	// BookTypeTxt txt格式
	BookTypeTxt BookType = "txt"
	// BookTypeEpub epub格式
	BookTypeEpub BookType = "epub"
)

// ReqBody 请求包
type ReqBody struct {
	// ArticaleURL 文章的链接地址, 系列文章中任意一篇即可
	// 例如： https://www.literotica.com/s/home-for-horny-monsters-ch-93
	ArticaleURL string
	// BookType 导出的电子书类型
	BookType BookType
}

// BasicInfo 基础信息
type BasicInfo struct {
	// Title 标题
	Title string
	// IndexURL 索引URL
	IndexURL []string
}

// SectionInfo 章节信息
type SectionInfo struct {
	// Title 章节标题
	Title string
	// Content 章节内容
	Content string
}

// ParseReqBody 解析包格式
func ParseReqBody() ReqBody {
	return ReqBody{
		ArticaleURL: url,
		BookType:    BookType(bookType),
	}
}
