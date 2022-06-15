package contracts_test

import (
	"os"
	"testing"

	"github.com/spacetab-io/mails-go/contracts"
	"github.com/spacetab-io/mails-go/errors"
	"github.com/stretchr/testify/assert"
)

var testAtt = contracts.Attachment{
	MimeType:     "text/plain; charset=utf-8",
	AttachMethod: contracts.AttachMethodFile,
	Filename:     "test.file",
	Name:         "test",
	Extension:    "file",
	Content:      []byte("some content"),
}

func TestNewAttachmentFromFile(t *testing.T) {
	type testCase struct {
		name string
		in   string
		exp  contracts.Attachment
		err  error
	}

	tcs := []testCase{
		{
			name: "normal attachment",
			in:   "./test.file",
			exp: contracts.Attachment{
				MimeType:     "text/plain; charset=utf-8",
				AttachMethod: contracts.AttachMethodFile,
				Name:         "test",
				Filename:     "test.file",
				Extension:    "file",
				Content:      []byte("some content"),
			},
			err: nil,
		},
		{
			name: "not existing file attachment",
			in:   "./not_existing_test.file",
			exp:  contracts.Attachment{},
			err:  os.ErrNotExist,
		},
		{
			name: "attachment is a dir",
			in:   "../contracts",
			exp:  contracts.Attachment{},
			err:  errors.ErrAttachmentIsNotAFile,
		},
		//{
		//	name: "not available to read file",
		//	in:   "./test_no_access.file",
		//	exp:  contracts.Attachment{},
		//	err:  os.ErrPermission,
		// },
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			att, err := contracts.NewAttachmentFromFile(tc.in)
			if tc.err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp, att)
		})
	}
}

func TestAttachment_IsEmpty(t *testing.T) {
	type testCase struct {
		name string
		in   contracts.Attachment
		exp  bool
	}

	tcs := []testCase{
		{
			name: "empty attachment",
			in:   contracts.Attachment{},
			exp:  true,
		},
		{
			name: "filled attachment",
			in:   testAtt,
			exp:  false,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.exp, tc.in.IsEmpty())
		})
	}
}

func TestAttachment_GetFileName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "test.file", testAtt.GetFileName())
}

func TestAttachment_GetMimeType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "text/plain; charset=utf-8", testAtt.GetMimeType())
}

func TestAttachment_GetName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "test", testAtt.GetName())
}

func TestAttachment_GetAttachMethod(t *testing.T) {
	t.Parallel()

	assert.Equal(t, contracts.AttachMethodFile, testAtt.GetAttachMethod())
}

func TestAttachment_GetContent(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []byte("some content"), testAtt.GetContent())
}
