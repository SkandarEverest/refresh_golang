package dto

import "mime/multipart"

type FileUploadRequest struct {
	FileHeader *multipart.FileHeader `form:"file" binding:"required"`
}

type FileUploadResponse struct {
	FileUrl string `json:"fileUrl"`
}
