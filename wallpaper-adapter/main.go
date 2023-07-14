package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

var (
	device   string // 设备，ipad 尺寸为 1536x2048
	filePath string

	ipadHeight = 2048
	ipadWidth  = 1536
)

func init() {
	flag.StringVar(&device, "device", "", "设备，目前只支持ipad")
	flag.StringVar(&filePath, "file", "", "文件路径，可以传图片，或者目录，传目录则处理整个目录下的文件")
	flag.Parse()
}

func main() {
	images := getImages(filePath)
	for _, f := range images {
		adaptOne(f)
	}
}

// 获取目录下的所有文件
func getImages(dirPath string) []string {
	// 获取目录下的所有文件和子目录
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		if isImageFile(filePath) {
			return []string{dirPath}
		}
		fmt.Println("无法读取目录：", err)
		return nil
	}

	var list []string
	// 遍历文件和子目录
	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		// 判断是否为目录
		if file.IsDir() {
			fmt.Println("子目录：", filePath)
			continue
		}
		// 判断是否为图片文件（可根据实际需要修改判断条件）
		if isImageFile(filePath) {
			fmt.Println("图片文件：", filePath)
			list = append(list, filePath)
		}
	}
	return list
}

// 判断是否为图片文件
func isImageFile(filePath string) bool {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	// 判断是否为图片扩展名（可根据实际需要添加更多判断条件）
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func adaptOne(filePath string) {
	// 读图
	img := read(filePath)

	// 如果是竖着的图
	// 把图片resize为  height=1536 width=2048
	if isTall(img) {
		img = resizeImage(img, ipadWidth, ipadHeight)
	}

	// 如果是横着的图
	// 把图片resize为 height=2048 width=1536
	img = resizeImage(img, ipadHeight, ipadWidth)

	// 保存修改后的图片
	if err := imaging.Save(img, "resized-"+filepath.Base(filePath)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("resize ok. ", filePath)
}

func read(filePath string) image.Image {
	// 打开原始图片文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 解码图片文件
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// isTall 图片是不是高的
func isTall(img image.Image) bool {
	bounds := img.Bounds()
	return bounds.Dy() > bounds.Dx()
}

// resizeImage 指定目标宽度和高度
func resizeImage(img image.Image, targetHeight, targetWidth int) image.Image {
	// 计算缩放后的宽度和高度
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	var newWidth, newHeight int

	// 计算缩放后的尺寸
	if srcWidth > srcHeight {
		newWidth = targetWidth
		newHeight = int(float64(srcHeight) * (float64(targetWidth) / float64(srcWidth)))
	} else {
		newWidth = int(float64(srcWidth) * (float64(targetHeight) / float64(srcHeight)))
		newHeight = targetHeight
	}

	// 进行缩放
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// 创建新的空白画布
	canvas := imaging.New(targetWidth, targetHeight, color.White)

	// 计算居中位置
	offsetX := (targetWidth - newWidth) / 2
	offsetY := (targetHeight - newHeight) / 2

	// 在画布上居中绘制缩放后的图片
	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+newWidth, offsetY+newHeight), resizedImg, image.Point{}, draw.Src)
	return canvas
}
