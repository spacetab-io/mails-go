package providers

import (
	"context"
	"crypto/tls"
	"fmt"

	cfgstructs "github.com/spacetab-io/configuration-structs-go/v2/contracts"
	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	customMime "github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/contracts"
	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTP struct {
	client      *mail.SMTPClient
	providerCfg mailing.MailProviderConfigInterface
}

func NewSMTP(providerCfg mailing.MailProviderConfigInterface) (SMTP, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return SMTP{}, fmt.Errorf("smtp provider config validation error: %w", err)
	}

	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = providerCfg.GetHostPort().GetHost()
	server.Port = int(providerCfg.GetHostPort().GetPort())
	server.Username = providerCfg.GetUsername()
	server.Password = providerCfg.GetPassword()
	server.Encryption = toProviderEncryption(providerCfg.GetEncryption())

	// Since v2.3.0 you can specified authentication type:
	// - PLAIN (default)
	// - LOGIN
	// - CRAM-MD5
	// - None
	server.Authentication = toProviderAuthType(providerCfg.GetAuthType())

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = providerCfg.GetConnectionTimeout()

	// Timeout for send the data and wait respond
	server.SendTimeout = providerCfg.GetSendTimeout()

	if providerCfg.GetEncryption() == mailing.MailProviderEncryptionNone {
		// Set TLSConfig to provide custom TLS configuration. For example,
		// to skip TLS verification (useful for testing):
		server.TLSConfig = &tls.Config{InsecureSkipVerify: true} //nolint: gosec
	}

	// SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		return SMTP{}, fmt.Errorf("smtp server connection error: %w", err)
	}

	return SMTP{client: smtpClient, providerCfg: providerCfg}, nil
}

func (o SMTP) Name() mailing.MailProviderName {
	return "smtp"
}

func (o SMTP) Send(_ context.Context, msg contracts.MessageInterface) error {
	// New email simple html with inline and CC
	email := mail.NewMSG().SetSubject(msg.GetSubject()).SetFrom(msg.GetFrom().String())

	tos := make([]string, 0, len(msg.GetTo().GetList()))
	for _, to := range msg.GetTo().GetList() {
		tos = append(tos, to.String())
	}

	email.AddTo(tos...)

	if !msg.GetCc().IsEmpty() {
		ccs := make([]string, 0, len(msg.GetCc().GetList()))
		for _, cc := range msg.GetCc().GetList() {
			ccs = append(ccs, cc.String())
		}

		email.AddCc(ccs...)
	}

	if !msg.GetBcc().IsEmpty() {
		bccs := make([]string, 0, len(msg.GetBcc().GetList()))
		for _, bcc := range msg.GetBcc().GetList() {
			bccs = append(bccs, bcc.String())
		}

		email.AddBcc(bccs...)
	}

	mime := mail.TextPlain

	if msg.GetMimeType() == customMime.TextHTML {
		mime = mail.TextHTML
	}

	email.SetBody(mime, string(msg.GetBody()))

	// also you can add body from []byte with SetBodyData, example:
	// email.SetBodyData(mail.TextHTML, []byte(htmlBody))
	// or alternative part
	// email.AddAlternativeData(mail.TextHTML, []byte(htmlBody))

	// add inline
	email.Attach(&mail.File{FilePath: "/path/to/image.png", Name: "Gopher.png", Inline: true})

	// you can add dkim signature to the email.
	// to add dkim, you need a private key already created one.
	if o.providerCfg.GetDKIMPrivateKey() != nil {
		options := dkim.NewSigOptions()
		k := o.providerCfg.GetDKIMPrivateKey()
		options.PrivateKey = []byte(*k)
		options.Domain = msg.GetFrom().GetDomain()
		options.Selector = "default"
		options.SignatureExpireIn = 3600
		options.Headers = []string{"from", "date", "mime-version", "received", "received"}
		options.AddSignatureTimestamp = true
		options.Canonicalization = "relaxed/relaxed"

		email.SetDkim(options)
	}

	// Call Send and pass the client
	if err := email.Send(o.client); err != nil {
		return fmt.Errorf("smtp email send error: %w", err)
	}

	// always check error after send
	if email.Error != nil {
		return email.Error
	}

	return nil
}

func toProviderAuthType(pat cfgstructs.AuthType) (at mail.AuthType) {
	switch pat {
	case cfgstructs.AuthTypePlain:
		at = mail.AuthPlain
	case cfgstructs.AuthTypeLogin:
		at = mail.AuthLogin
	case cfgstructs.AuthTypeCRAMMD5:
		at = mail.AuthCRAMMD5
	case cfgstructs.AuthTypeNone:
		at = mail.AuthNone
	case cfgstructs.AuthTypeBasic, cfgstructs.AuthTypeJWT:
		at = mail.AuthPlain
	default:
		at = mail.AuthPlain
	}

	return at
}

func toProviderEncryption(enc mailing.MailProviderEncryption) (me mail.Encryption) {
	switch enc {
	case mailing.MailProviderEncryptionNone:
		me = mail.EncryptionNone
	case mailing.MailProviderEncryptionSSL:
		me = mail.EncryptionSSLTLS
	case mailing.MailProviderEncryptionTLS:
		me = mail.EncryptionTLS
	case mailing.MailProviderEncryptionSSLTLS:
		me = mail.EncryptionSSLTLS
	case mailing.MailProviderEncryptionSTARTTLS:
		me = mail.EncryptionSTARTTLS
	default:
		me = mail.EncryptionNone
	}

	return
}
