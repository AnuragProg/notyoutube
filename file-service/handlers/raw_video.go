package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/anuragprog/notyoutube/file-service/configs"
	databaseRepo "github.com/anuragprog/notyoutube/file-service/repository/database"
	mqRepo "github.com/anuragprog/notyoutube/file-service/repository/mq"
	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
	errType "github.com/anuragprog/notyoutube/file-service/types/errors"
	mqType "github.com/anuragprog/notyoutube/file-service/types/mq"
)

func PostRawVideoHandler(db databaseRepo.Database, store *storeRepo.StoreManager, mq *mqRepo.MessageQueueManager) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		if len(files) != 1 {
			return errType.NewAPIError(http.StatusBadRequest, "exactly 1 file required")
		}
		file := files[0]

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
		var assumedUploadSpeedBytesPerSec float64 = 1 << 20 // kb per second
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


		// push event to kafka queue
		// TODO: add logging mechanism for kafka events 
		go func(metadata databaseType.RawVideoMetadata){
			if err := mq.PublishToRawVideoTopic(mqType.FromRawVideoMetadataToProtoRawVideoMetadata(metadata)); err != nil {
				fmt.Println(err.Error())
			}
		}(metadata)

		return c.JSON(http.StatusCreated, metadata)
	}
}

func GetRawVideoMetadatasHandler(db databaseRepo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < configs.DEFAULT_PAGE_START {
			page = configs.DEFAULT_PAGE_START
		}
		pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
		if err != nil || pageSize < configs.DEFAULT_PAGE_SIZE {
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
		return c.JSON(http.StatusOK, metadatas)
	}
}

func GetRawVideoHandler(db databaseRepo.Database, store *storeRepo.StoreManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		videoId := c.QueryParam("video_id")
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

		var assumedDownloadSpeedBytesPerSec float64 = 1 << 20 // kb per second
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
		if _, err := io.Copy(c.Response().Writer, file); err != nil {
			return errType.IntoAPIError(err, http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}
