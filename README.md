# mails-go

Golang library for email sending.

## Providers

List of available providers:

* [Sendgrid](github.com/sendgrid/sendgrid-go)
* [Mandrill](github.com/mattbaird/gochimp)
* [Mailgun](github.com/mailgun/mailgun-go/v4)
* [SMTP](github.com/xhit/go-simple-mail/v2)
* log
* file

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/spacetab-io/configuration-structs-go/v2/configuration/mimetype"
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/mails-go/contracts"
	"github.com/spacetab-io/mails-go/providers"
)

func main() {
	// 1. Get Provider config (with should implement mailing.MailProviderConfigInterface
	// For Example, Sendgrid Config
	sendgridCfg := mailing.SendgridConfig{
		Enabled:     true,
		Key:         "APIKey",
		SendTimeout: 5 * time.Second,
	}

	// 2. Initiate provider
	// Sendgrid provider
	sendgrid, err := providers.NewSendgrid(sendgridCfg)
	if err != nil {
		panic(err)
	}

	// 3. Prepare Message with should implement contracts.MessageInterface
	msg := contracts.Message{
		To: mailing.MailAddressList{
			{Email: "toOne@spacetab.io", Name: "To One"},
			{Email: "totwo@spacetab.io", Name: "To Two"},
		},
		MimeType: mime.TextPlain,
		Subject:  "Test email",
		Content:  []byte("test email content"),
	}

	// 4. Send message
	if err := sendgrid.Send(context.Background(), msg); err != nil {
		panic(err)
	}
}
```