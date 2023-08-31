// Package processor 处理整个过程
package processor

import "github.com/r0oot/tools/internal/literotica/protocol"

// NextSectionGetter 取下一章节
type NextSectionGetter func() (*protocol.SectionInfo, error)

// Reader 读接口
type Reader interface {
	// GetInfoAndGetter 取基础信息, 返回Book的基础信息和第一章节的getter
	GetInfoAndGetter(articleURL string) (*protocol.BasicInfo, NextSectionGetter, error)
}

// Writer 写接口
type Writer interface {
	// Create 创建文件
	Create(basicInfo *protocol.BasicInfo) error
	// AddSection 添加章节
	AddSection(sectionInfo *protocol.SectionInfo) error
	// Save 保存
	Save() error
}

// Processor 处理器
type Processor struct {
	r Reader
	w Writer
}

// New 创建实例
func New(r Reader, w Writer) *Processor {
	return &Processor{
		r: r,
		w: w,
	}
}

func (p *Processor) Do(reqBody *protocol.ReqBody) error {
	basicInfo, getter, err := p.r.GetInfoAndGetter(reqBody.ArticaleURL)
	if err != nil {
		return err
	}
	if err := p.w.Create(basicInfo); err != nil {
		return err
	}
	for {
		sectionInfo, err := getter()
		if err != nil {
			return err
		}
		if sectionInfo == nil {
			break
		}
		if err := p.w.AddSection(sectionInfo); err != nil {
			return err
		}
	}
	return p.w.Save()
}
