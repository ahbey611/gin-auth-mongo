package databases

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

var MinioClient *minio.Client
var defaultBucketName string

func InitMinio() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	endpoint := os.Getenv("MINIO_URL")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_ACCESS_KEY_SECRET")
	useSSL := false // TODO: change to true when deploy to production

	// Initialize minio client object.
	myminioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	MinioClient = myminioClient
	defaultPublicBucketName := os.Getenv("MINIO_PUBLIC_BUCKET_NAME")
	defaultPrivateBucketName := os.Getenv("MINIO_PRIVATE_BUCKET_NAME")

	// check if bucket exists
	publicBucketExists, err := MinioClient.BucketExists(context.Background(), defaultPublicBucketName)
	if err != nil {
		panic(err)
	}
	if !publicBucketExists {
		err = MinioClient.MakeBucket(context.Background(), defaultPublicBucketName, minio.MakeBucketOptions{})
		if err != nil {
			panic(err)
		}
		log.Println("Bucket created:", defaultPublicBucketName)

		// MinioClient.SetBucketPolicy(context.Background(), defaultPublicBucketName)
	} else {
		log.Println("Bucket already exists:", defaultPublicBucketName)
	}

	privateBucketExists, err := MinioClient.BucketExists(context.Background(), defaultPrivateBucketName)
	if err != nil {
		log.Fatalln(err)
	}
	if !privateBucketExists {
		MinioClient.MakeBucket(context.Background(), defaultPrivateBucketName, minio.MakeBucketOptions{})
	} else {
		log.Println("Bucket already exists:", defaultPrivateBucketName)
	}
}

func PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (err error) {
	if bucketName == "" {
		bucketName = defaultBucketName
	}

	_, err = MinioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		log.Println(err)
	}
	return err
}

func TagObject(ctx context.Context, bucketName string, objectName string, tag string) error {
	tags, err := tags.NewTags(map[string]string{
		"tag": tag,
	}, false)
	if err != nil {
		log.Println(err)
	}
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	err = MinioClient.PutObjectTagging(ctx, bucketName, objectName, tags, minio.PutObjectTaggingOptions{})
	if err != nil {
		log.Println(err)
	}
	return err
}

func GetObjectTag(ctx context.Context, bucketName string, objectName string) (map[string]string, error) {
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	result, err := MinioClient.GetObjectTagging(ctx, bucketName, objectName, minio.GetObjectTaggingOptions{})
	// log.Println("result: ", result)
	if err != nil {
		log.Println(err)
	}
	return result.ToMap(), err
}

func GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (obj *minio.Object, err error) {
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	obj, err = MinioClient.GetObject(ctx, bucketName, objectName, opts)
	if err != nil {
		log.Println(err)
	}
	return obj, err
}

func RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) (err error) {
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	err = MinioClient.RemoveObject(ctx, bucketName, objectName, opts)
	if err != nil {
		log.Println(err)
	}
	return err
}

func StatObject(ctx context.Context, bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	objectInfo, err := MinioClient.StatObject(ctx, bucketName, objectName, opts)
	if err != nil {
		log.Println(err)
	}
	return objectInfo, err
}
