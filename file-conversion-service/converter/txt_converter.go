package converter

import (
	"bytes"
	"os"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type TxtConverter struct{}

func NewTxtConverter() *TxtConverter {
	return &TxtConverter{}
}

func (c *TxtConverter) ConvertToPDF(content []byte) ([]byte, error) {

	if _, err := os.Stat("./assets/fonts/DejaVuSans.ttf"); os.IsNotExist(err) {
		return nil, err
	}

	decoder := unicode.UTF8.NewDecoder()
	utf8Content, _, err := transform.Bytes(decoder, content)
	if err != nil {
		return nil, err
	}
	text := string(utf8Content)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("DejaVu", "", "./assets/fonts/DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 12)
	pdf.MultiCell(0, 10, text, "", "", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	return buf.Bytes(), err
}
