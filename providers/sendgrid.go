package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/contracts"
)

type Sendgrid struct {
	client      *sendgrid.Client
	providerCfg mailing.MailProviderConfigInterface
}

func NewSendgrid(providerCfg mailing.MailProviderConfigInterface) (Sendgrid, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return Sendgrid{}, fmt.Errorf("sendgrid provider config validation error: %w", err)
	}

	return Sendgrid{client: sendgrid.NewSendClient(providerCfg.GetPassword()), providerCfg: providerCfg}, nil
}

func (o Sendgrid) Name() mailing.MailProviderName {
	return "sendgridAPI"
}

func (o Sendgrid) Send(ctx context.Context, msg contracts.MessageInterface) error {
	var content *mail.Content

	switch msg.GetMimeType() {
	case mime.TextHTML, mime.TextPlain:
		content = mail.NewContent(msg.GetMimeType().String(), string(msg.GetBody()))
	default:
		content = mail.NewContent(mime.TextPlain.String(), string(msg.GetBody()))
	}

	message := mail.NewV3Mail()

	if !msg.GetFrom().IsEmpty() {
		message.SetFrom(mail.NewEmail(msg.GetFrom().GetName(), msg.GetFrom().GetEmail()))
	}

	message.AddPersonalizations(o.getPersonalization(msg))
	message.AddContent(content)

	if !msg.GetReplyTo().IsEmpty() {
		message.ReplyTo = mail.NewEmail(msg.GetReplyTo().GetName(), msg.GetReplyTo().GetEmail())
	}

	message.Subject = msg.GetSubject()

	ctx, cancel := context.WithTimeout(ctx, o.providerCfg.GetSendTimeout())

	defer cancel()

	response, err := o.client.SendWithContext(ctx, message)
	if err == nil && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("sendgrid send message error: %d %s", response.StatusCode, response.Body) //nolint: goerr113
	} else if err != nil {
		return fmt.Errorf("sendgrid email send error: %w", err)
	}

	return nil
}

func (o Sendgrid) getPersonalization(msg contracts.MessageInterface) *mail.Personalization {
	p := mail.NewPersonalization()

	if !msg.GetFrom().IsEmpty() {
		p.AddFrom(mail.NewEmail(msg.GetFrom().GetName(), msg.GetFrom().GetEmail()))
	}

	tos := make([]*mail.Email, 0)
	for _, to := range msg.GetTo().GetList() {
		tos = append(tos, mail.NewEmail(to.GetName(), to.GetEmail()))
	}

	p.AddTos(tos...)

	if !msg.GetCc().IsEmpty() {
		ccs := make([]*mail.Email, 0, len(msg.GetCc().GetList()))
		for _, cc := range msg.GetCc().GetList() {
			ccs = append(ccs, mail.NewEmail(cc.GetName(), cc.GetEmail()))
		}

		p.AddCCs(ccs...)
	}

	if !msg.GetBcc().IsEmpty() {
		bccs := make([]*mail.Email, 0, len(msg.GetBcc().GetList()))
		for _, bcc := range msg.GetBcc().GetList() {
			bccs = append(bccs, mail.NewEmail(bcc.GetName(), bcc.GetEmail()))
		}

		p.AddBCCs(bccs...)
	}

	p.Subject = msg.GetSubject()

	return p
}
