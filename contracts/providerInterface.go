package contracts

import (
	"context"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
)

type ProviderInterface interface {
	Name() mailing.MailProviderName
	Send(ctx context.Context, msg MessageInterface) error
}
