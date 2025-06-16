package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RtfConverter struct{}

func NewRtfConverter() *RtfConverter {
	return &RtfConverter{}
}

func (c *RtfConverter) ConvertToPDF(rtfData []byte) ([]byte, error) {
	tmpDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("rtf_conv_%d", time.Now().UnixNano()),
	)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.rtf")
	if err := os.WriteFile(inputPath, rtfData, 0644); err != nil {
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
