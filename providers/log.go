package providers

import (
	"context"
	"fmt"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/mails-go/contracts"
)

type LogProvider struct {
	providerCfg mailing.MailProviderConfigInterface
	logger      contracts.LoggerInterface
}

func NewLogProvider(providerCfg mailing.MailProviderConfigInterface, logger contracts.LoggerInterface) (LogProvider, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return LogProvider{}, fmt.Errorf("log provider config validation error: %w", err)
	}

	return LogProvider{providerCfg: providerCfg, logger: logger}, nil
}

func (o LogProvider) Name() mailing.MailProviderName {
	return mailing.MailProviderLogs
}

func (o LogProvider) Send(_ context.Context, msg contracts.MessageInterface) error {
	o.logger.Printf("email by [%s]:\n%s", o.Name(), msg.String())

	return nil
}
