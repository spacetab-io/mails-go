package providers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spacetab-io/configuration-structs-go/v2/mailing"
	"github.com/spacetab-io/mails-go/contracts"
)

type File struct {
	providerCfg mailing.MailProviderConfigInterface
}

func NewFileProvider(providerCfg mailing.MailProviderConfigInterface) (File, error) {
	if _, err := providerCfg.Validate(); err != nil {
		return File{}, fmt.Errorf("file provider config validation error: %w", err)
	}

	if !fileExists(providerCfg.GetHostPort().GetHost()) {
		if err := createFile(providerCfg.GetHostPort().GetHost()); err != nil {
			return File{}, fmt.Errorf("email file create error: %w", err)
		}
	} else {
		if err := os.Chtimes(providerCfg.GetHostPort().GetHost(), time.Now().Local(), time.Now().Local()); err != nil {
			return File{}, fmt.Errorf("email file touch error: %w", err)
		}
	}

	return File{providerCfg: providerCfg}, nil
}

func (f File) Name() mailing.MailProviderName {
	return mailing.MailProviderFile
}

func (f File) Send(_ context.Context, msg contracts.MessageInterface) error {
	file, err := os.OpenFile(
		f.providerCfg.GetHostPort().GetHost(),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0o600, // nolint: gomnd
	)
	if err != nil {
		return fmt.Errorf("file open on email send error: %w", err)
	}

	defer file.Close()

	if _, err = file.WriteString(msg.String() + "\n"); err != nil {
		return fmt.Errorf("file write on email send error: %w", err)
	}

	return nil
}

func createFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("mailing file %s create error: %w", filePath, err)
	}

	defer file.Close()

	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
