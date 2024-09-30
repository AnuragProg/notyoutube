package handlers

import (
	"context"
	"errors"
	"io"
	"math"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/anuragprog/notyoutube/file-service/configs"
	databaseRepo "github.com/anuragprog/notyoutube/file-service/repository/database"
	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
	errType "github.com/anuragprog/notyoutube/file-service/types/errors"
)

func PostRawVideoHandler(db databaseRepo.Database, store *storeRepo.StoreManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return errType.IntoAPIError(err, http.StatusBadRequest, "post video request must be a multipart form data and contain file with 'file' as key and actual file as value")
		}

		files := form.File["file"]
		if len(files) == 0 {
			return errType.NewAPIError(http.StatusBadRequest, "no video files found")
		}
		if len(files) != 1 {
			return errType.NewAPIError(http.StatusNotImplemented, "currently we accept only one file, and will in future have functionality to serve multiple files")
		}

		generatedMetadatas := make([]databaseType.RawVideoMetadata, 0, len(files))

		for _, file := range files {
			contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
			metadata := databaseType.RawVideoMetadata{
				Filename:    file.Filename,
				ContentType: contentType,
				FileSize:    file.Size,
				CreatedAt:   time.Now().UTC(),
			}

			// create database entry
			ctx, cancel := context.WithTimeout(context.Background(), configs.DEFAULT_TIMEOUT)
			defer cancel()
			generatedMetadata, err := db.CreateRawVideoMetadata(ctx, metadata)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return errType.IntoAPIError(err, http.StatusServiceUnavailable, "database call timed out")
			case err != nil:
				return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
			}

			// create object in store
			var assumedUploadSpeedBytesPerSec float64 = 1 << 20 // bytes per second
			var bufferTimeSecs int64 = 10                       // seconds
			expectedUploadTimeSecs := int64(math.Ceil((float64(file.Size) / assumedUploadSpeedBytesPerSec))) + bufferTimeSecs
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(expectedUploadTimeSecs))
			defer cancel()
			fileReader, err := file.Open()
			if err != nil {
				return errType.IntoAPIError(err, http.StatusBadRequest, err.Error())
			}
			err = store.Upload(ctx, storeRepo.RAW_VIDEO, generatedMetadata.Id, fileReader, file.Size, generatedMetadata.ContentType)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return errType.IntoAPIError(err, http.StatusServiceUnavailable, "storage call timed out")
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

func GetRawVideoMetadatasHandler(db databaseRepo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page")
		pageSize := c.QueryInt("page_size")
		if page < configs.DEFAULT_PAGE_START {
			page = configs.DEFAULT_PAGE_START
		}
		if pageSize < configs.DEFAULT_PAGE_SIZE {
			pageSize = configs.DEFAULT_PAGE_SIZE
		}
		ctx, cancel := context.WithTimeout(context.Background(), configs.DEFAULT_TIMEOUT)
		defer cancel()
		metadatas, err := db.ListRawVideosMetadata(ctx, page, pageSize)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return errType.IntoAPIError(err, http.StatusServiceUnavailable, "database call timed out")
		case err != nil:
			return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
		}
		return c.Status(http.StatusOK).JSON(metadatas)
	}
}

func GetRawVideoHandler(db databaseRepo.Database, store *storeRepo.StoreManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		videoId := c.Params("video_id")
		ctx, cancel := context.WithTimeout(context.Background(), configs.DEFAULT_TIMEOUT)
		defer cancel()
		metadata, err := db.GetRawVideoMetadata(ctx, videoId)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return errType.IntoAPIError(err, http.StatusServiceUnavailable, "database call timed out")
		case errors.Is(err, errType.InvalidQuery):
			return errType.IntoAPIError(err, http.StatusBadRequest, err.Error())
		case errors.Is(err, errType.RecordNotFound):
			return errType.IntoAPIError(err, http.StatusNotFound, err.Error())
		case err != nil:
			return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
		}

		var assumedDownloadSpeedBytesPerSec float64 = 1 << 20 // bytes per second
		var bufferTimeSecs int64 = 10                         // seconds
		expectedDownloadTimeSecs := int64(math.Ceil((float64(metadata.FileSize) / assumedDownloadSpeedBytesPerSec))) + bufferTimeSecs
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(expectedDownloadTimeSecs))
		defer cancel()
		file, err := store.Download(ctx, storeRepo.RAW_VIDEO, metadata.Id)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return errType.IntoAPIError(err, http.StatusServiceUnavailable, "storage call timed out")
		case err != nil:
			return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
		}
		defer file.Close()

		c.Set("Content-Type", metadata.ContentType)
		if _, err := io.Copy(c.Response().BodyWriter(), file); err != nil {
			return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}
