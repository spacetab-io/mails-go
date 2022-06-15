package providers

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/contracts"
)

type Mailgun struct {
	client      *mailgun.MailgunImpl
	providerCfg mailing.MailProviderConfigInterface
}

func NewMailgun(providerCfg mailing.MailProviderConfigInterface) (Mailgun, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return Mailgun{}, fmt.Errorf("mailgun provider config validation error: %w", err)
	}

	mg := mailgun.NewMailgun(providerCfg.GetUsername(), providerCfg.GetPassword())
	mg.SetAPIBase(providerCfg.GetHostPort().String())

	return Mailgun{client: mg, providerCfg: providerCfg}, nil
}

func (o Mailgun) Name() mailing.MailProviderName {
	return "mailgunAPI"
}

func (o Mailgun) Send(ctx context.Context, msg contracts.MessageInterface) error {
	tos := make([]string, 0)
	for _, to := range msg.GetTo().GetList() {
		tos = append(tos, to.String())
	}

	message := o.client.NewMessage(msg.GetFrom().String(), msg.GetSubject(), string(msg.GetBody()), tos...)

	if !msg.GetCc().IsEmpty() {
		for _, cc := range msg.GetCc().GetList() {
			message.AddCC(cc.String())
		}
	}

	if !msg.GetBcc().IsEmpty() {
		for _, bcc := range msg.GetBcc().GetList() {
			message.AddBCC(bcc.String())
		}
	}

	if !msg.GetReplyTo().IsEmpty() {
		message.SetReplyTo(msg.GetReplyTo().String())
	}

	if msg.GetMimeType() == mime.TextHTML {
		message.SetHtml(string(msg.GetBody()))
	}

	if o.providerCfg.GetDKIMPrivateKey() != nil {
		message.SetDKIM(true)
	}

	ctx, cancel := context.WithTimeout(ctx, o.providerCfg.GetSendTimeout())

	defer cancel()

	if _, _, err := o.client.Send(ctx, message); err != nil {
		return fmt.Errorf("%s send message error: %w", o.Name(), err)
	}

	return nil
}
