package handlers

import (
	"context"
	"errors"
	"math"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"

	databaseRepo "github.com/anuragprog/notyoutube/file-service/repository/database"
	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
	errType "github.com/anuragprog/notyoutube/file-service/types/errors"
)

func PostRawVideoHandler(
	db databaseRepo.Database,
	store storeRepo.StoreManager ,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return errType.IntoAPIError(err, http.StatusBadRequest, "Post video request must be a multipart form data and contain file with 'file' as key and actual file as value")
		}

		files := form.File["file"]
		if len(files) == 0 {
			return errType.NewAPIError(http.StatusBadRequest, "no video files found")
		}
		if len(files) != 1 {
			return errType.NewAPIError(http.StatusNotImplemented, "currently we accept only one file, and will in future have functionality to serve multiple files")
		}

		
		generatedMetadatas := make([]databaseType.RawVideoMetadata, len(files))

		for _, file := range files {
			contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
			metadata := databaseType.RawVideoMetadata{
				Filename: file.Filename,
				ContentType: contentType,
				FileSize: file.Size,
				CreatedAt: time.Now().UTC(),
			}

			// create database entry
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			generatedMetadata, err := db.CreateRawVideoMetadata(ctx, metadata)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return errType.IntoAPIError(err, http.StatusInternalServerError, "Database call timed out")
			case err != nil:
				return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
			}

			// create object in store
			var assumedUploadSpeedBytesPerSec float64 = 1<<20 // bytes per second
			var bufferTimeSecs int64 = 10 // seconds
			expectedUploadTimeSecs := int64(math.Ceil((float64(file.Size)/assumedUploadSpeedBytesPerSec))) + bufferTimeSecs
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(expectedUploadTimeSecs))
			defer cancel()
			fileReader, err := file.Open()
			if err != nil {
				return errType.IntoAPIError(err, http.StatusBadRequest, err.Error())
			}
			err = store.Upload(ctx, storeRepo.RAW_VIDEO, generatedMetadata.Id, fileReader, file.Size, generatedMetadata.ContentType)
			switch{
			case errors.Is(err, context.DeadlineExceeded):
				return errType.IntoAPIError(err, http.StatusInternalServerError, "Storage call timed out")
			case err != nil:
				return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
			}
			fileReader.Close()

			// append 
			generatedMetadatas = append(generatedMetadatas, generatedMetadata)
		}
		return c.Status(http.StatusCreated).JSON(generatedMetadatas)
	}
}
