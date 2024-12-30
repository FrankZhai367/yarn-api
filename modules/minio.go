package modules

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/cenkalti/backoff"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/yarn/backend/config"
)

type MinioServer struct {
	conf   *config.Config
	Client *minio.Client
}

func NewMinioServer(conf *config.Config) (*MinioServer, error) {
	minioClient, err := CreateMinioClientOrFatal(&conf.Minio)
	if err != nil {
		logrus.Errorln("create minio client err: ", err)
		return nil, err
	}
	server := &MinioServer{
		conf:   conf,
		Client: minioClient,
	}
	err = server.createBucketName()
	if err != nil {
		logrus.Errorln("create minio bucketName err:", err)
		return nil, err
	}
	logrus.Infoln("create minio bucketName success")
	err = server.createMinioPolicy()
	if err != nil {
		logrus.Errorln("create minio policy err:", err)
		return nil, err
	}
	logrus.Infoln("create minio policy success")
	return server, nil
}

func createCredentialProviderChain(accessKey, secretKey string) *credentials.Credentials {
	if accessKey != "" && secretKey != "" {
		return credentials.NewStaticV4(accessKey, secretKey, "")
	}
	providers := []credentials.Provider{
		&credentials.EnvMinio{},
		&credentials.EnvAWS{},
		&credentials.IAM{
			Client: &http.Client{
				Transport: http.DefaultTransport,
			},
		},
	}
	return credentials.New(&credentials.Chain{Providers: providers})
}

func CreateMinioClient(conf *config.Minio) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	cred := createCredentialProviderChain(conf.AccessKey, conf.SecretAccessKey)
	opts := &minio.Options{
		Creds:  cred,
		Secure: conf.Secure,
	}
	minioClient, err := minio.New(endpoint, opts)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while creating minio client: %+v", err))
	}
	return minioClient, nil
}

func CreateMinioClientOrFatal(conf *config.Minio) (*minio.Client, error) {
	var minioClient *minio.Client
	var err error
	var operation = func() error {
		minioClient, err = CreateMinioClient(conf)
		if err != nil {
			return err
		}
		return nil
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = conf.ConnectionTimeout * time.Second
	err = backoff.Retry(operation, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return minioClient, nil
}

func (server *MinioServer) createBucketName() error {
	exists, err := server.Client.BucketExists(context.Background(), server.conf.Minio.BucketName)
	if exists {
		return nil
	}
	if err != nil {
		return err
	}
	err = server.Client.MakeBucket(context.Background(), server.conf.Minio.BucketName, minio.MakeBucketOptions{
		Region:        server.conf.Minio.Region,
		ObjectLocking: true,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (server *MinioServer) createMinioPolicy() error {
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, server.conf.Minio.BucketName)
	//fmt.Println(policy)
	return server.Client.SetBucketPolicy(context.Background(), server.conf.Minio.BucketName, policy)
}

func (server *MinioServer) AddLocalfile(localFile, minioPath, contentType string) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	destPath := path.Join(server.conf.Minio.BasePath, minioPath)
	_, err := server.Client.FPutObject(context.Background(), server.conf.Minio.BucketName, destPath, localFile,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to store %v: %v", minioPath, err.Error()))
		return errors.Errorf("Failed to store %v: %v", minioPath, err.Error())
	}
	logrus.Infoln(fmt.Sprintf("Succeed to store %v", minioPath))
	return nil
}

func (server *MinioServer) DeleteFile(filePath string) error {
	err := server.Client.RemoveObject(context.Background(), server.conf.Minio.BucketName, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to delete %v: %v", filePath, err.Error()))
		return errors.Errorf("Failed to delete %v: %v", filePath, err.Error())
	}
	logrus.Infoln(fmt.Sprintf("Succeed to delete %v", filePath))
	return nil
}

func (server *MinioServer) GetGileToLocal(filePath, localFilePath string) error {
	err := server.Client.FGetObject(context.Background(), server.conf.Minio.BucketName, filePath, localFilePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (server *MinioServer) GetFile(filePath string) ([]byte, error) {
	reader, err := server.Client.GetObject(context.Background(), server.conf.Minio.BucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Errorf("Failed to get %v: %v", filePath, err.Error())
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return nil, errors.Errorf("Failed to get %v: %v", filePath, err.Error())
	}

	tmpBytes := buf.Bytes()
	if server.conf.Minio.Multipart {
		re := regexp.MustCompile(`\w+;chunk-signature=\w+`)
		tmpBytes = []byte(re.ReplaceAllString(string(tmpBytes), ""))
	}
	return tmpBytes, nil
}

func (server *MinioServer) GetObject(fileName string) (*minio.Object, error) {
	object, err := server.Client.GetObject(context.Background(), server.conf.Minio.BucketName, fileName, minio.GetObjectOptions{})
	return object, err
}
