package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/jung-kurt/gofpdf"
)

type DicomConverter struct{}

func NewDicomConverter() *DicomConverter {
	return &DicomConverter{}
}

func (c *DicomConverter) ConvertToPDF(dcmData []byte) ([]byte, error) {
	tmpDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("dcm_%d", time.Now().UnixNano()),
	)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.dcm")
	if err := os.WriteFile(inputPath, dcmData, 0644); err != nil {
		return nil, fmt.Errorf("failed to save DICOM file: %w", err)
	}

	pngPath := filepath.Join(tmpDir, "output.png")
	cmd := exec.Command(
		"dcm2img",
		inputPath,
		pngPath,
		"--write-png",
		"--quiet",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("dcm2img failed: %s\n%w", string(output), err)
	}

	img, err := imaging.Open(pngPath)
	if err != nil {
		return nil, fmt.Errorf("PNG decode failed: %w", err)
	}

	var jpgBuf bytes.Buffer
	if err := imaging.Encode(&jpgBuf, img, imaging.JPEG, imaging.JPEGQuality(90)); err != nil {
		return nil, fmt.Errorf("JPEG encode failed: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.RegisterImageOptionsReader("image.jpg",
		gofpdf.ImageOptions{ImageType: "JPG"},
		bytes.NewReader(jpgBuf.Bytes()))

	pdf.Image("image.jpg", 10, 10, 190, 0, false, "", 0, "")

	var pdfBuf bytes.Buffer
	if err := pdf.Output(&pdfBuf); err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}

	return pdfBuf.Bytes(), nil
}
