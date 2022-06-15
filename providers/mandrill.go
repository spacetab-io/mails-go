package providers

import (
	"context"
	"fmt"

	"github.com/mattbaird/gochimp"
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/contracts"
)

type Mandrill struct {
	mandrillAPI *gochimp.MandrillAPI
	providerCfg mailing.MailProviderConfigInterface
}

func NewMandrill(providerCfg mailing.MailProviderConfigInterface) (Mandrill, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return Mandrill{}, fmt.Errorf("mandrill provider config validation error: %w", err)
	}

	api, err := gochimp.NewMandrill(providerCfg.GetPassword())
	if err != nil {
		return Mandrill{}, fmt.Errorf("mandrill client init error: %w", err)
	}

	if providerCfg.GetSendTimeout() != 0 {
		api.Timeout = providerCfg.GetSendTimeout()
	}

	return Mandrill{mandrillAPI: api, providerCfg: providerCfg}, nil
}

func (o Mandrill) Name() mailing.MailProviderName {
	return "mandrillAPI"
}

func (o Mandrill) Send(_ context.Context, msg contracts.MessageInterface) error {
	tos := make([]gochimp.Recipient, 0)

	for _, to := range msg.GetTo().GetList() {
		tos = append(tos, gochimp.Recipient{
			Name:  to.GetName(),
			Email: to.GetEmail(),
		})
	}

	if !msg.GetCc().IsEmpty() {
		for _, cc := range msg.GetCc().GetList() {
			tos = append(tos, gochimp.Recipient{
				Name:  cc.GetName(),
				Email: cc.GetEmail(),
				Type:  "cc",
			})
		}
	}

	if !msg.GetBcc().IsEmpty() {
		for _, bcc := range msg.GetBcc().GetList() {
			tos = append(tos, gochimp.Recipient{
				Name:  bcc.GetName(),
				Email: bcc.GetEmail(),
				Type:  "bcc",
			})
		}
	}

	message := gochimp.Message{
		Subject:   msg.GetSubject(),
		FromName:  msg.GetFrom().GetName(),
		FromEmail: msg.GetFrom().GetEmail(),
		To:        tos,
	}

	switch msg.GetMimeType() {
	case mime.TextHTML:
		message.Html = string(msg.GetBody())
	case mime.TextPlain:
		message.Text = string(msg.GetBody())
	default:
		message.Text = string(msg.GetBody())
	}

	if !msg.GetReplyTo().IsEmpty() {
		message.Headers = map[string]string{
			"Reply-To": msg.GetReplyTo().String(),
		}
	}

	if _, err := o.mandrillAPI.MessageSend(message, o.providerCfg.IsAsync()); err != nil {
		return fmt.Errorf("mandrill email send error: %w", err)
	}

	return nil
}
