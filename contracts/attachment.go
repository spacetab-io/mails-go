package contracts

import (
	"os"
	"path/filepath"
	"strings"

	mime "github.com/gabriel-vasile/mimetype"
	"github.com/spacetab-io/mails-go/errors"
)

type AttachMethod string

const (
	AttachMethodInline AttachMethod = "inline"
	AttachMethodFile   AttachMethod = "file"
)

type Attachment struct {
	MimeType     string
	AttachMethod AttachMethod
	Filename     string
	Name         string
	Extension    string
	Content      []byte
}

func NewAttachmentFromFile(filePath string) (a Attachment, err error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}

	if fi.IsDir() {
		return a, errors.ErrAttachmentIsNotAFile
	}

	a.Content, err = os.ReadFile(filePath)
	if err != nil {
		return
	}

	fileName := fi.Name()
	m, _ := mime.DetectFile(filePath)

	a.MimeType = m.String()
	a.Extension = strings.Trim(filepath.Ext(filePath), ".")
	a.Name = fileName[:len(fileName)-len(filepath.Ext(fileName))]
	a.Filename = fi.Name()
	a.AttachMethod = AttachMethodFile

	return a, nil
}

func (a Attachment) IsEmpty() bool {
	return len(a.Content) == 0
}

func (a Attachment) GetFileName() string {
	return a.Filename
}

func (a Attachment) GetMimeType() string {
	return a.MimeType
}

func (a Attachment) GetContent() []byte {
	return a.Content
}

func (a Attachment) GetName() string {
	return a.Name
}

func (a Attachment) GetAttachMethod() AttachMethod {
	return a.AttachMethod
}
