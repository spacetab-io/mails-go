package mails_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go"
	"github.com/spacetab-io/mails-go/contracts"
	"github.com/spacetab-io/mails-go/providers"
	"github.com/stretchr/testify/assert"
)

func TestMailing_Send(t *testing.T) {
	type inStruct struct {
		mailing func(writer io.Writer) mails.Mailing
		msg     contracts.Message
	}
	type testCase struct {
		name string
		in   inStruct
		exp  func() string
		err  error
	}
	basicMsg := contracts.Message{
		To:          mailing.MailAddressList{mailing.MailAddress{Email: "toOne@spacetab.io", Name: "To One"}, mailing.MailAddress{Email: "totwo@spacetab.io", Name: "To Two"}},
		MimeType:    mime.TextPlain,
		Subject:     "Test email",
		Content:     []byte("test email content"),
		Attachments: nil,
	}

	fullMsg := basicMsg
	fullMsg.From = mailing.MailAddress{Email: "from@spacetab.io", Name: "FromName"}
	fullMsg.ReplyTo = mailing.MailAddress{Email: "replyTo@spacetab.io", Name: "replyToName"}
	fullMsg.Cc = mailing.MailAddressList{mailing.MailAddress{Email: "cc@spacetab.io", Name: "Carbon Copy"}}
	fullMsg.Bcc = mailing.MailAddressList{mailing.MailAddress{Email: "bcc@spacetab.io", Name: "Blind Carbon Copy"}}

	basicMailCfg := mailing.MessagingConfig{
		From:    mailing.MailAddress{Email: "robot@spacetab.io", Name: "Spacetab Robot"},
		ReplyTo: mailing.MailAddress{Email: "feedback@spacetab.io", Name: "Spacetab Feedback"},
	}

	fullMailCfg := basicMailCfg
	fullMailCfg.SubjectPrefix = "[test]"

	prefix := "email by [logs]:\n"

	tcs := []testCase{
		{
			name: "basic email config and full message",
			in: inStruct{
				msg: fullMsg,
				mailing: func(writer io.Writer) mails.Mailing {
					mockLogger := mails.NewLogger(writer)
					mockProvider, _ := providers.NewLogProvider(mailing.LogsConfig{}, mockLogger)

					return mails.NewMailingForProvider(mockProvider, basicMailCfg)
				},
			},
			exp: func() string {
				return prefix + fullMsg.String()
			},
		},
		{
			name: "full email config and basic message",
			in: inStruct{msg: basicMsg, mailing: func(writer io.Writer) mails.Mailing {
				mockLogger := mails.NewLogger(writer)
				mockProvider, _ := providers.NewLogProvider(mailing.LogsConfig{}, mockLogger)

				return mails.NewMailingForProvider(mockProvider, fullMailCfg)
			}},
			exp: func() string {
				msg := basicMsg
				msg.From = fullMailCfg.From
				msg.Subject = fullMailCfg.SubjectPrefix + " " + msg.Subject
				msg.ReplyTo = fullMailCfg.ReplyTo

				return prefix + msg.String()
			},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bb := &bytes.Buffer{}
			m := tc.in.mailing(bb)
			ctx := context.Background()
			err := m.Send(ctx, &tc.in.msg)
			if tc.err != nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp(), bb.String())
		})
	}
}
