package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BmpConverter struct{}

func NewBmpConverter() *BmpConverter {
	return &BmpConverter{}
}

func (c *BmpConverter) ConvertToPDF(bmpData []byte) ([]byte, error) {
	tmpDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("bmp_conv_%d", time.Now().UnixNano()),
	)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.bmp")
	if err := os.WriteFile(inputPath, bmpData, 0644); err != nil {
		return nil, fmt.Errorf("failed to save input file: %w", err)
	}

	cmd := exec.Command(
		"C:\\Program Files\\LibreOffice\\program\\soffice.exe",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", tmpDir,
		inputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("libreoffice error: %s\n%w", string(output), err)
	}

	outputPath := filepath.Join(tmpDir, "input.pdf")
	return os.ReadFile(outputPath)
}
