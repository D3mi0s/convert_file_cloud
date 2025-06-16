package converter

import (
	"fmt"
)

type Converter interface {
	ConvertToPDF(input []byte) ([]byte, error)
}

func NewConverter() *FileConverter {
	return &FileConverter{}
}

type FileConverter struct {
}

func (c *FileConverter) ConvertFile(fileData []byte, mimeType string) ([]byte, error) {
	switch {

	case mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		mimeType == "application/msword":
		return NewDocxConverter().ConvertToPDF(fileData)

	case mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		mimeType == "application/vnd.ms-excel":
		return NewXlsxConverter().ConvertToPDF(fileData)

	case mimeType == "text/plain":
		return NewTxtConverter().ConvertToPDF(fileData)

	case mimeType == "image/png",
		mimeType == "image/jpeg":
		return NewImageConverter().ConvertToPDF(fileData)
	case mimeType == "application/vnd.openxmlformats-officedocument.presentationml.presentation", // PPTX
		mimeType == "application/vnd.ms-powerpoint": // PPT
		return NewPptConverter().ConvertToPDF(fileData)
	case mimeType == "text/rtf",
		mimeType == "application/rtf":
		return NewRtfConverter().ConvertToPDF(fileData)
	case mimeType == "image/bmp",
		mimeType == "image/x-bmp",
		mimeType == "image/x-ms-bmp":
		return NewBmpConverter().ConvertToPDF(fileData)
	case mimeType == "application/dicom",
		mimeType == "application/octet-stream" && isDicomFile(fileData):
		return NewDicomConverter().ConvertToPDF(fileData)

	default:
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}
}

func isDicomFile(data []byte) bool {

	if len(data) < 132 {
		return false
	}
	return string(data[128:132]) == "DICM"
}
