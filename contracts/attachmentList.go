package contracts

type MessageAttachmentList []Attachment

func (mal MessageAttachmentList) GetFileNames() []string {
	if mal.IsEmpty() {
		return nil
	}

	fileNames := make([]string, 0, len(mal.GetList()))

	for _, file := range mal.GetList() {
		fileNames = append(fileNames, file.GetFileName())
	}

	return fileNames
}

func (mal MessageAttachmentList) GetList() []MessageAttachmentInterface {
	mm := make([]MessageAttachmentInterface, 0, len(mal))

	for _, ma := range mal {
		mm = append(mm, ma)
	}

	return mm
}

func (mal MessageAttachmentList) IsEmpty() bool {
	return len(mal) == 0
}
