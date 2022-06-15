package mails

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/mails-go/contracts"
	"github.com/spacetab-io/mails-go/providers"
)

type Mailing struct {
	provider contracts.ProviderInterface
	msgCfg   mailing.MessagingConfigInterface
}

func NewMailing(providerCfg mailing.MailProviderConfigInterface, msgCfg mailing.MessagingConfigInterface) (Mailing, error) {
	var (
		provider contracts.ProviderInterface
		err      error
	)

	switch providerCfg.Name() {
	case mailing.MailProviderLogs:
		var w io.Writer

		switch providerCfg.GetHostPort().GetHost() {
		case "stdout":
			w = os.Stdout
		case "stderr":
			w = os.Stderr
		default:
			w = os.Stdout
		}

		logger := NewLogger(w)
		provider, err = providers.NewLogProvider(mailing.LogsConfig{}, logger)
	case mailing.MailProviderFile:
		provider, err = providers.NewFileProvider(providerCfg)
	case mailing.MailProviderMailgun:
		provider, err = providers.NewMailgun(providerCfg)
	case mailing.MailProviderMandrill:
		provider, err = providers.NewMandrill(providerCfg)
	case mailing.MailProviderSendgrid:
		provider, err = providers.NewSendgrid(providerCfg)
	case mailing.MailProviderSMTP:
		provider, err = providers.NewSMTP(providerCfg)
	default:
		return Mailing{}, mailing.ErrUnknownProvider
	}

	if err != nil {
		return Mailing{}, fmt.Errorf("provider init error: %w", err)
	}

	return Mailing{provider: provider, msgCfg: msgCfg}, nil
}

func NewMailingForProvider(provider contracts.ProviderInterface, msgCfg mailing.MessagingConfigInterface) Mailing {
	return Mailing{provider: provider, msgCfg: msgCfg}
}

func (m Mailing) Send(ctx context.Context, msg contracts.MessageInterface) error {
	if msg.GetFrom().IsEmpty() && !m.msgCfg.GetFrom().IsEmpty() {
		_ = msg.SetFrom(m.msgCfg.GetFrom())
	}

	if !m.msgCfg.GetReplyTo().IsEmpty() && msg.GetReplyTo().IsEmpty() {
		_ = msg.SetReplyTo(m.msgCfg.GetReplyTo())
	}

	if m.msgCfg.GetSubjectPrefix() != "" {
		_ = msg.SetSubject(fmt.Sprintf(
			"%s %s",
			strings.TrimSpace(m.msgCfg.GetSubjectPrefix()),
			strings.TrimSpace(msg.GetSubject()),
		))
	}

	if err := m.provider.Send(ctx, msg); err != nil {
		return fmt.Errorf("mailing send error: %w", err)
	}

	return nil
}
