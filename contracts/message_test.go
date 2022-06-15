package contracts_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/configuration-structs-go/v2/mime"
	"github.com/spacetab-io/mails-go/contracts"
	"github.com/spacetab-io/mails-go/errors"
	"github.com/stretchr/testify/assert"
)

func TestMessage_SetFrom(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddress
		exp  string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddress{Name: "from", Email: "from@email.com"},
			exp:  "\"from\" <from@email.com>",
		},
		{
			name: "empty address",
			in:   mailing.MailAddress{},
			exp:  "",
			err:  errors.ErrEmptyAddress,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetFrom(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetFrom().String())
		})
	}
}

func TestMessage_SetReplyTo(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddress
		exp  string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddress{Name: "replyTo", Email: "replyTo@email.com"},
			exp:  "\"replyTo\" <replyTo@email.com>",
		},
		{
			name: "empty address",
			in:   mailing.MailAddress{},
			exp:  "",
			err:  errors.ErrEmptyAddress,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetReplyTo(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetReplyTo().String())
		})
	}
}

func TestMessage_SetTo(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddress
		exp  []string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddress{Name: "to", Email: "to@email.com"},
			exp:  []string{"\"to\" <to@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddress{},
			exp:  nil,
			err:  errors.ErrEmptyAddress,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetTo(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetTo().GetStringList())
		})
	}
}

func TestMessage_SetTos(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddressList
		exp  []string
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddressList{mailing.MailAddress{Name: "toOne", Email: "toOne@email.com"}, mailing.MailAddress{Name: "toTwo", Email: "toTwo@email.com"}},
			exp:  []string{"\"toOne\" <toOne@email.com>", "\"toTwo\" <toTwo@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddressList{},
			exp:  nil,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			msg.SetTos(tc.in)

			assert.Equal(t, tc.exp, msg.GetTo().GetStringList())
		})
	}
}

func TestMessage_SetCc(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddress
		exp  []string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddress{Name: "cc", Email: "cc@email.com"},
			exp:  []string{"\"cc\" <cc@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddress{},
			exp:  nil,
			err:  errors.ErrEmptyAddress,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetCc(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetCc().GetStringList())
		})
	}
}

func TestMessage_SetCcs(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddressList
		exp  []string
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddressList{mailing.MailAddress{Name: "ccOne", Email: "ccOne@email.com"}, mailing.MailAddress{Name: "ccTwo", Email: "ccTwo@email.com"}},
			exp:  []string{"\"ccOne\" <ccOne@email.com>", "\"ccTwo\" <ccTwo@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddressList{},
			exp:  nil,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			msg.SetCcs(tc.in)

			assert.Equal(t, tc.exp, msg.GetCc().GetStringList())
		})
	}
}

func TestMessage_SetBcc(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddress
		exp  []string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddress{Name: "bcc", Email: "bcc@email.com"},
			exp:  []string{"\"bcc\" <bcc@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddress{},
			exp:  nil,
			err:  errors.ErrEmptyAddress,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetBcc(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, tc.err, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetBcc().GetStringList())
		})
	}
}

func TestMessage_SetBccs(t *testing.T) {
	type testCase struct {
		name string
		in   mailing.MailAddressList
		exp  []string
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   mailing.MailAddressList{mailing.MailAddress{Name: "bccOne", Email: "bccOne@email.com"}, mailing.MailAddress{Name: "bccTwo", Email: "bccTwo@email.com"}},
			exp:  []string{"\"bccOne\" <bccOne@email.com>", "\"bccTwo\" <bccTwo@email.com>"},
		},
		{
			name: "empty address",
			in:   mailing.MailAddressList{},
			exp:  nil,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			msg.SetBccs(tc.in)

			assert.Equal(t, tc.exp, msg.GetBcc().GetStringList())
		})
	}
}

func TestMessage_SetSubject(t *testing.T) {
	type testCase struct {
		name string
		in   string
		exp  string
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   "subject",
			exp:  "subject",
		},
		{
			name: "empty data",
			in:   "",
			exp:  "",
			err:  errors.ErrEmptyData,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetSubject(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetSubject())
		})
	}
}

func TestMessage_SetHTML(t *testing.T) {
	type expStruct struct {
		mime    mime.Type
		content []byte
	}
	type testCase struct {
		name string
		in   []byte
		exp  expStruct
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   []byte("<body>text<body>"),
			exp: expStruct{
				mime:    mime.TextHTML,
				content: []byte("<body>text<body>"),
			},
		},
		{
			name: "empty data",
			in:   nil,
			exp:  expStruct{},
			err:  errors.ErrEmptyData,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetHTML(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp.content, msg.GetBody())
			assert.Equal(t, tc.exp.mime, msg.GetMimeType())
		})
	}
}

func TestMessage_SetPlainText(t *testing.T) {
	type expStruct struct {
		mime    mime.Type
		content []byte
	}
	type testCase struct {
		name string
		in   []byte
		exp  expStruct
		err  error
	}

	tcs := []testCase{
		{
			name: "correct setting",
			in:   []byte("plain text"),
			exp: expStruct{
				mime:    mime.TextPlain,
				content: []byte("plain text"),
			},
		},
		{
			name: "empty data",
			in:   nil,
			exp:  expStruct{},
			err:  errors.ErrEmptyData,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := contracts.Message{}

			err := msg.SetPlainText(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp.content, msg.GetBody())
			assert.Equal(t, tc.exp.mime, msg.GetMimeType())
		})
	}
}

func TestMessage_AddAttachment(t *testing.T) {
	type inStruct struct {
		filePath string
	}
	type testCase struct {
		name string
		in   inStruct
		exp  []string
		err  error
	}

	tcs := []testCase{
		{
			name: "filled attachment",
			in:   inStruct{filePath: "test.file"},
			exp:  []string{"test.file"},
			err:  nil,
		},
		{
			name: "empty attachment",
			in:   inStruct{},
			exp:  nil,
			err:  errors.ErrEmptyData,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				msg contracts.Message
				att contracts.Attachment
				err error
			)

			if tc.in.filePath != "" {
				att, err = contracts.NewAttachmentFromFile(tc.in.filePath)
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			err = msg.AddAttachment(att)
			if tc.err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, msg.GetAttachments().GetFileNames())
		})
	}
}

func TestMessage_String(t *testing.T) {
	t.Parallel()

	format := `=======================
from: %s
to: %s
cc: %s
bc: %s
replyTo: %s
subject: %s

body:
%s
=======================
`

	msg := contracts.Message{
		To:          mailing.MailAddressList{mailing.MailAddress{Email: "toOne@spacetab.io", Name: "To One"}, mailing.MailAddress{Email: "totwo@spacetab.io", Name: "To Two"}},
		MimeType:    mime.TextPlain,
		Subject:     "Test email",
		Content:     []byte("test email content"),
		Attachments: nil,
	}

	expString := fmt.Sprintf(format,
		msg.GetFrom().String(),
		strings.Join(msg.GetTo().GetStringList(), ", "),
		strings.Join(msg.GetCc().GetStringList(), ", "),
		strings.Join(msg.GetBcc().GetStringList(), ", "),
		msg.GetReplyTo().String(),
		msg.GetSubject(),
		string(msg.GetBody()),
	)

	assert.Equal(t, expString, msg.String())
}
