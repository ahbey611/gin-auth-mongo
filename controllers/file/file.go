package file

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"

	"gin-auth-mongo/databases"
	"net/http"
	"time"

	"gin-auth-mongo/utils/file"
	"gin-auth-mongo/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// 上传文件
func UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	folder := "avatar"

	// 生成文件名
	objectName := folder + "/" + header.Filename

	// 设置 MinIO 上传选项
	opts := minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")}

	// 上传文件到 MinIO
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	objectSize := header.Size
	if err := databases.PutObject(ctx, "", objectName, file, objectSize, opts); err != nil {
		response.BadRequestWithMessage(c, "Failed to upload file")
		return
	}

	response.SuccessWithData(c, gin.H{"filename": objectName})
}

// 下载文件
func DownloadFile(c *gin.Context) {
	// 从 URL 查询参数中获取文件名
	filename := c.Query("filename")
	if filename == "" {
		response.BadRequestWithMessage(c, "Filename is required")
		return
	}

	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	bucketName := os.Getenv("MINIO_PUBLIC_BUCKET_NAME")

	// 从 MinIO 获取文件对象
	obj, err := databases.MinioClient.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		response.BadRequestWithMessage(c, "Failed to retrieve file from MinIO")
		return
	}
	defer obj.Close()

	// 获取文件的信息以设置内容类型和内容长度
	stat, err := obj.Stat()
	if err != nil {
		response.BadRequestWithMessage(c, "Failed to retrieve file information")
		return
	}

	// 设置响应头，确保浏览器可以下载文件
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))

	// 将文件内容写入响应中
	if _, err := io.Copy(c.Writer, obj); err != nil {
		response.BadRequestWithMessage(c, "Failed to download file")
		return
	}
}

// generate presigned url
func GeneratePresignedUrl(c *gin.Context) {
	filename := c.Query("filename")
	expires := c.DefaultQuery("expires", "10") // 默认 10 分钟

	// 将 expires 转换为整数分钟数
	expiresMinutes, err := strconv.Atoi(expires)
	if err != nil {
		response.BadRequestWithMessage(c, "Invalid expires parameter")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	bucketName := os.Getenv("MINIO_PUBLIC_BUCKET_NAME")
	presignedUrl, err := databases.MinioClient.PresignedGetObject(ctx, bucketName, filename, time.Duration(expiresMinutes)*time.Minute, nil)
	if err != nil {
		response.BadRequestWithMessage(c, "Failed to generate presigned url")
		return
	}

	response.SuccessWithData(c, gin.H{"presignedUrl": presignedUrl.String()})
}

func DecryptFileName(c *gin.Context) {
	fileName := c.PostForm("filename")
	if fileName == "" {
		response.BadRequestWithMessage(c, "Filename is required")
		return
	}

	fileName = strings.Split(fileName, ".")[0]
	if fileName == "" {
		response.BadRequestWithMessage(c, "Invalid filename")
		return
	}

	fileInfo, err := file.DecryptFileName(fileName)
	if err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return
	}
	response.SuccessWithData(c, fileInfo)
}
