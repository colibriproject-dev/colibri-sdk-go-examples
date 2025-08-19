package repositories

import (
	"mime/multipart"
	"os"
	"testing"

	"github.com/colibriproject-dev/colibri-sdk-go-examples/school-module/src/infra/storages"
	"github.com/colibriproject-dev/colibri-sdk-go/pkg/base/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var documentStorage = storages.NewDocumentS3Storage()

func TestDocumentStorage_Upload(t *testing.T) {
	t.Run("Should upload document to storage", func(t *testing.T) {
		var file multipart.File
		file, err := os.Open(test.MountAbsolutPath("../../../development-environment/files/img.png"))
		assert.NoError(t, err)

		result, err := documentStorage.Upload(ctx, uuid.New(), &file)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
