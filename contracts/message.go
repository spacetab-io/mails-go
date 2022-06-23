package contracts

import (
	"fmt"
	"strings"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/errors"
)

type Message struct {
	From    mailing.MailAddress
	ReplyTo mailing.MailAddress
	To      mailing.MailAddressList
	Cc      mailing.MailAddressList
	Bcc     mailing.MailAddressList

	MimeType mime.Type
	Subject  string
	Content  []byte

	Attachments MessageAttachmentList
}

func (mm *Message) SetFrom(addr mailing.MailAddressInterface) error {
	if addr.IsEmpty() {
		return errors.ErrEmptyAddress
	}

	mm.From = mailing.NewMailAddressFromInterface(addr)

	return nil
}

func (mm *Message) SetTo(addr mailing.MailAddressInterface) error {
	if addr.IsEmpty() {
		return errors.ErrEmptyAddress
	}

	mm.To = append(mm.To, mailing.NewMailAddressFromInterface(addr))

	return nil
}

func (mm *Message) SetTos(addrs mailing.MailAddressListInterface) {
	for _, addr := range addrs.GetList() {
		mm.To = append(mm.To, mailing.NewMailAddressFromInterface(addr))
	}
}

func (mm *Message) SetCc(addr mailing.MailAddressInterface) error {
	if addr.IsEmpty() {
		return errors.ErrEmptyAddress
	}

	mm.Cc = append(mm.Cc, mailing.NewMailAddressFromInterface(addr))

	return nil
}

func (mm *Message) SetCcs(addrs mailing.MailAddressListInterface) {
	for _, addr := range addrs.GetList() {
		mm.Cc = append(mm.Cc, mailing.NewMailAddressFromInterface(addr))
	}
}

func (mm *Message) SetBcc(addr mailing.MailAddressInterface) error {
	if addr.IsEmpty() {
		return errors.ErrEmptyAddress
	}

	mm.Bcc = append(mm.Bcc, mailing.NewMailAddressFromInterface(addr))

	return nil
}

func (mm *Message) SetBccs(addrs mailing.MailAddressListInterface) {
	for _, addr := range addrs.GetList() {
		mm.Bcc = append(mm.Bcc, mailing.NewMailAddressFromInterface(addr))
	}
}

func (mm *Message) SetReplyTo(addr mailing.MailAddressInterface) error {
	if addr.IsEmpty() {
		return errors.ErrEmptyAddress
	}

	mm.ReplyTo = mailing.NewMailAddressFromInterface(addr)

	return nil
}

func (mm *Message) SetSubject(sbj string) error {
	if sbj == "" {
		return fmt.Errorf("%w: %s", errors.ErrEmptyData, "subject")
	}

	mm.Subject = sbj

	return nil
}

func (mm *Message) SetMimeType(typ mime.Type) {
	mm.MimeType = typ
}

func (mm *Message) SetHTML(msg []byte) error {
	if msg == nil {
		mm.emptyContent()

		return fmt.Errorf("%w: %s", errors.ErrEmptyData, "content")
	}

	mm.MimeType = mime.TextHTML
	mm.Content = msg

	return nil
}

func (mm *Message) SetPlainText(msg []byte) error {
	if msg == nil {
		mm.emptyContent()

		return fmt.Errorf("%w: %s", errors.ErrEmptyData, "content")
	}

	mm.MimeType = mime.TextPlain
	mm.Content = msg

	return nil
}

func (mm *Message) AddAttachment(file MessageAttachmentInterface) error {
	if file.IsEmpty() {
		return fmt.Errorf("%w: %s", errors.ErrEmptyData, "attachment")
	}

	mm.Attachments = append(mm.Attachments, Attachment{
		MimeType:     file.GetMimeType(),
		AttachMethod: file.GetAttachMethod(),
		Name:         file.GetName(),
		Filename:     file.GetFileName(),
		Content:      file.GetContent(),
	})

	return nil
}

func (mm *Message) AddAttachments(files ...MessageAttachmentInterface) error {
	for _, file := range files {
		if err := mm.AddAttachment(file); err != nil {
			return err
		}
	}

	return nil
}

func (mm Message) GetAttachments() MessageAttachmentListInterface {
	return mm.Attachments
}

func (mm Message) GetFrom() mailing.MailAddressInterface {
	return mm.From
}

func (mm Message) GetTo() mailing.MailAddressListInterface {
	return mm.To
}

func (mm Message) GetCc() mailing.MailAddressListInterface {
	return mm.Cc
}

func (mm Message) GetBcc() mailing.MailAddressListInterface {
	return mm.Bcc
}

func (mm Message) GetReplyTo() mailing.MailAddressInterface {
	return mm.ReplyTo
}

func (mm Message) GetMimeType() mime.Type {
	return mm.MimeType
}

func (mm Message) GetBody() []byte {
	return mm.Content
}

func (mm Message) GetSubject() string {
	return mm.Subject
}

func (mm Message) String() string {
	msgFormat := `=======================
from: %s
to: %s
cc: %s
bc: %s
replyTo: %s
subject: %s

body:
%s
=======================
`

	return fmt.Sprintf(msgFormat,
		mm.GetFrom().String(),
		strings.Join(mm.GetTo().GetStringList(), ", "),
		strings.Join(mm.GetCc().GetStringList(), ", "),
		strings.Join(mm.GetBcc().GetStringList(), ", "),
		mm.GetReplyTo().String(),
		mm.GetSubject(),
		string(mm.GetBody()),
	)
}

func (mm *Message) emptyContent() {
	mm.MimeType = ""
	mm.Content = nil
}
