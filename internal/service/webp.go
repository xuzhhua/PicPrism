package service

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

const jpegQuality = 85

// generateOptimized 将图片编码为 JPEG 写入 path（质量 85）
func generateOptimized(img image.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_ = jpeg.Encode(f, img, &jpeg.Options{Quality: jpegQuality})
}

// generateThumb 生成宽度不超过 maxWidth 的 JPEG 缩略图
func generateThumb(img image.Image, path string, maxWidth int) {
	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	newW, newH := origW, origH
	if origW > maxWidth {
		newW = maxWidth
		newH = origH * maxWidth / origW
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	// 白色背景（处理透明图）
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			dst.Set(x, y, color.White)
		}
	}
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_ = jpeg.Encode(f, dst, &jpeg.Options{Quality: jpegQuality})
}
