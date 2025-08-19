//go:generate mockgen -source document_storage.go -destination mock/document_storage_mock.go -package storagesmock
package storages

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/colibriproject-dev/colibri-sdk-go/pkg/storage"
	"github.com/google/uuid"
)

type IDocumentStorage interface {
	Upload(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error)
}

type DocumentS3Storage struct {
	bucket string
}

func NewDocumentS3Storage() *DocumentS3Storage {
	return &DocumentS3Storage{
		bucket: "meu-bucket",
	}
}

func (s *DocumentS3Storage) Upload(ctx context.Context, id uuid.UUID, file *multipart.File) (string, error) {
	return storage.UploadFile(ctx, s.bucket, fmt.Sprintf("STUDENT-DOCUMENT-%s", id.String()), file)
}
