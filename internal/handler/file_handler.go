package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/SkandarEverest/refresh-golang/internal/dto"
	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/SkandarEverest/refresh-golang/internal/usecase"
	"github.com/SkandarEverest/refresh-golang/pkg/helper"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type FileHandler struct {
	Log     *logrus.Logger
	UseCase *usecase.FileUseCase
}

func NewFileHandler(useCase *usecase.FileUseCase, logger *logrus.Logger) *FileHandler {
	return &FileHandler{Log: logger,
		UseCase: useCase}
}

func (c *FileHandler) UploadFile(ctx echo.Context) error {

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.ServerError(err.Error()))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.ServerError(err.Error()))
	}

	if isValid := c.isValidFile(fileHeader); !isValid {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("file is not valid"))
	}

	fileUrl, err := c.UseCase.UploadFile(file)
	if err != nil {
		ex, ok := err.(*exception.CustomError)
		if ok {
			c.Log.Warnf("Failed to Upload File : %+v", err)
			return ctx.JSON(ex.StatusCode, ex)
		}
		panic(err)
	}

	return ctx.JSON(http.StatusOK, &dto.FileUploadResponse{
		FileUrl: helper.DerefString(fileUrl, ""),
	})
}

func (c *FileHandler) isValidFile(fileHeader *multipart.FileHeader) bool {
	if fileHeader.Size > 2*1024*1024 || fileHeader.Size < 10*1024 {
		return false
	}

	fileType := fileHeader.Header.Get("Content-Type")

	return fileType == "image/jpeg" || fileType == "image/jpg" || fileType == "image/png"
}
