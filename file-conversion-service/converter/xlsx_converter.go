package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type XlsxConverter struct{}

func NewXlsxConverter() *XlsxConverter {
	return &XlsxConverter{}
}

func (c *XlsxConverter) ConvertToPDF(xlsxData []byte) ([]byte, error) {
	tmpDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("xlsx_conv_%d", time.Now().UnixNano()),
	)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.xlsx")
	if err := os.WriteFile(inputPath, xlsxData, 0644); err != nil {
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
