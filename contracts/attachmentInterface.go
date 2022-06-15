package contracts

type MessageAttachmentListInterface interface {
	GetList() []MessageAttachmentInterface
	IsEmpty() bool
	GetFileNames() []string
}

type MessageAttachmentInterface interface {
	IsEmpty() bool
	GetFileName() string
	GetName() string
	GetMimeType() string
	GetContent() []byte
	GetAttachMethod() AttachMethod
}
