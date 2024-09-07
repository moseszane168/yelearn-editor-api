/**
 * Minio工具类
 */

package minioutil

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var minioClient *minio.Client
var BucketName = "crf-mold"
var ctx = context.Background()

/**
 * 初始化Minio客户端连接
 */
func Init() {
	endpoint := fmt.Sprintf("%s:%s", viper.GetString("minio.host"), viper.GetString("minio.port"))
	useSSL := false

	// 初始化MinIO客户端对象
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("minio.accessKeyID"), viper.GetString("minio.secretAccessKey"), ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.WithError(err).Error("连接minio失败")
	}

	// 创建一个新的bucket
	err = minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	if err != nil {
		// 检查我们是否已经存在了这个bucket
		exists, errBucketExists := minioClient.BucketExists(ctx, BucketName)
		if errBucketExists != nil || !exists {
			logrus.WithError(err).Error("创建bucket失败")
		}
	} else {
		logrus.Info("Successfully created %s\n", BucketName)
	}

	logrus.Info("Minio初始化成功")
}

/**
 * 上传本地文件到Minio
 */
func UploadMinio(filePath, objectName string) minio.UploadInfo {
	// 上传文件
	info, err := minioClient.FPutObject(ctx, BucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return info
}
