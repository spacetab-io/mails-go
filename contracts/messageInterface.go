package contracts

import (
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
)

type MessageInterface interface {
	SetFrom(addr mailing.MailAddressInterface) error
	SetTo(addr mailing.MailAddressInterface) error
	SetTos(addrs mailing.MailAddressListInterface)
	SetCc(addr mailing.MailAddressInterface) error
	SetCcs(addrs mailing.MailAddressListInterface)
	SetBcc(addr mailing.MailAddressInterface) error
	SetBccs(addrs mailing.MailAddressListInterface)
	SetReplyTo(addr mailing.MailAddressInterface) error
	SetSubject(sbj string) error
	SetHTML(msg []byte) error
	SetPlainText(msg []byte) error
	AddAttachment(file MessageAttachmentInterface) error
	AddAttachments(files ...MessageAttachmentInterface) error

	GetFrom() mailing.MailAddressInterface
	GetTo() mailing.MailAddressListInterface
	GetCc() mailing.MailAddressListInterface
	GetBcc() mailing.MailAddressListInterface
	GetReplyTo() mailing.MailAddressInterface
	GetSubject() string
	GetMimeType() mime.Type
	GetBody() []byte
	GetAttachments() MessageAttachmentListInterface

	String() string
}
