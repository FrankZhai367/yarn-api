package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/yarn/backend/config"
	"github.com/yarn/backend/db"
	"github.com/yarn/backend/modules"
	"gorm.io/gorm"
)

type ImageService struct {
	db          *gorm.DB
	conf        *config.Config
	minioClient *modules.MinioServer
}

func NewImageService(db *gorm.DB, conf *config.Config) *ImageService {
	minioServer, _ := modules.NewMinioServer(conf)

	return &ImageService{
		db:          db,
		conf:        conf,
		minioClient: minioServer,
	}
}
func (service *ImageService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/image")

	r.POST("/upload", service.upload)
}

func (service *ImageService) GetImage(openid string) Resp {

	var image db.Image

	err := service.db.First(&image, "openid = ?", openid).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "用户奖励不存在"}
	}
	return Resp{Code: 0, Data: image, Msg: ""}
}

func (service *ImageService) upload(c *gin.Context) {
	file, err := c.FormFile("file")
	openid := c.GetHeader("openid")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}

	var bucketName = service.conf.Minio.BucketName
	var basePath = "upload"
	objectName := file.Filename
	filePath := fmt.Sprintf("%s/%s", basePath, file.Filename)

	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	defer fileHeader.Close()

	// 上传文件

	_, err = service.minioClient.Client.PutObject(c, bucketName, filePath, fileHeader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}

	service.db.Create(&db.Image{Filename: filePath, Openid: openid})

	presignedURL, err := service.minioClient.Client.PresignedGetObject(context.TODO(), service.conf.Minio.BucketName, filePath, 10*time.Minute, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}

	log.Printf("Successfully uploaded %s and presigned URL is %s\n", objectName, presignedURL)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "操作成功",
		"data": map[string]interface{}{
			"url":      fmt.Sprintf("%s/%s/%s", service.conf.Minio.AdvertisedHost, service.conf.Minio.BucketName, filePath),
			"filePath": filePath,
		},
	})
}
