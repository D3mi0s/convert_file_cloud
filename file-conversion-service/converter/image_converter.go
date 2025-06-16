package converter

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/jung-kurt/gofpdf"
)

type ImageConverter struct{}

func NewImageConverter() *ImageConverter {
	return &ImageConverter{}
}

func (c *ImageConverter) ConvertToPDF(imgData []byte) ([]byte, error) {

	contentType := http.DetectContentType(imgData)
	var imageType string
	switch contentType {
	case "image/png":
		imageType = "png"
	case "image/jpeg":
		imageType = "jpg"
	default:
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	imgW := float64(img.Bounds().Dx())
	imgH := float64(img.Bounds().Dy())
	pageW, pageH := 190.0, 277.0

	ratio := min(pageW/imgW, pageH/imgH)
	imgW *= ratio
	imgH *= ratio

	opts := gofpdf.ImageOptions{
		ImageType: imageType,
		ReadDpi:   true,
	}

	pdf.RegisterImageOptionsReader("image."+imageType, opts, bytes.NewReader(imgData))

	x := (210 - imgW) / 2
	y := (297 - imgH) / 2

	pdf.ImageOptions("image."+imageType, x, y, imgW, imgH, false, opts, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}

	return buf.Bytes(), nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
