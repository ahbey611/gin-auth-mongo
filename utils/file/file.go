package file

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"

	// "crypto/sha256"

	// "encoding/base64"
	// "encoding/hex"
	// "log"

	// "encoding/base32"
	// "encoding/base32"
	// "encoding/base64"
	"encoding/base64"
	// "encoding/hex"
	// "encoding/base64"
	// "encoding/base64"
	// "encoding/hex"
	// "encoding/hex"
	"errors"

	// "log"
	// "os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gin-auth-mongo/databases"
	"gin-auth-mongo/utils/consts"

	// "gin-auth-mongo/utils/datetime"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

var ALLOWED_IMAGE_EXTENSIONS = []string{".png", ".jpg", ".jpeg", ".svg", ".gif", ".webp"}
var ALLOWED_IMAGE_CONTENT_TYPES = []string{"image/png", "image/jpeg", "image/svg+xml", "image/gif", "image/webp"}

// get file name and extension, eg: avatar.png -> avatar, png
func GetFileNameAndExtension(fileName string) (string, string) {
	name := filepath.Base(fileName)
	// replace all '/' to '-' to prevent directory traversal
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "_", "-")
	if len(name) > 25 {
		name = name[:25]
	}
	extension := filepath.Ext(name)
	name = strings.TrimSuffix(name, extension)
	return name, extension
}

// CheckFileIsImage checks if the uploaded file is an image based on its MIME type.
func CheckFileIsImage(contentType string) bool {

	// Check if the content type matches any of the allowed image content types
	for _, allowedContentType := range ALLOWED_IMAGE_CONTENT_TYPES {
		if contentType == allowedContentType {
			return true
		}
	}
	return false
}

func CheckImageSize(size int64) bool {
	return size <= consts.MAX_IMAGE_FILE_SIZE
}

// generate file name max length is 255
func GenerateUniqueFileName(folder string, fileName string) string {
	// return folder + "/" + header.Filename
	name, extension := GetFileNameAndExtension(fileName)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	newFileName := name + "=" + uuid.New().String() + "$" + now
	length := len(folder + "/" + newFileName + extension)
	if length > 255 {
		newFileName = newFileName[:255-len(extension)-len(folder+"/")]
	}
	return folder + "/" + newFileName + extension
}

func GenerateEncryptedFileName(fileName string) (string, error) {

	name, extension := GetFileNameAndExtension(fileName)
	encryptionKey := consts.FILE_ENCRYPTION_KEY
	if encryptionKey == "" {
		return "", errors.New("encryption key is not set")
	}
	encryptedFileName, err := encrypt(name, encryptionKey)
	if err != nil {
		return "", err
	}
	return encryptedFileName + extension, nil
}

func DecryptFileName(encryptedFileName string) (map[string]string, error) {

	encryptionKey := consts.FILE_ENCRYPTION_KEY
	fileName, err := decrypt(encryptedFileName, encryptionKey)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"rawFileName": fileName,
	}, nil
}

// https://dev.to/breda/secret-key-encryption-with-go-using-aes-316d

func encrypt(data string, key string) (string, error) {

	aes, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.New("failed to create AES cipher")
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", errors.New("failed to create GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", errors.New("failed to generate nonce")
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(data), nil)
	// return base64.StdEncoding.EncodeToString(cipherText), nil
	// return base64.RawStdEncoding.EncodeToString(cipherText), nil
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
	// return hex.EncodeToString(cipherText), nil

}

func decrypt(data string, key string) (string, error) {
	aes, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.New("failed to create AES cipher")
	}

	// decodedData, err := base64.StdEncoding.DecodeString(data)
	// decodedData, err := base64.RawStdEncoding.DecodeString(data)
	decodedData, err := base64.RawURLEncoding.DecodeString(data)
	// decodedData, err := hex.DecodeString(data)
	if err != nil {
		log.Println("err: ", err)
		return "", errors.New("failed to decode data by base64")
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", errors.New("failed to create GCM")
	}

	nonceSize := gcm.NonceSize()

	nonce, cipherText := decodedData[:nonceSize], decodedData[nonceSize:]
	plainText, err := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return "", errors.New("failed to decrypt data")
	}
	return string(plainText), nil
}

func CheckFileIsExist(bucketName string, path string) bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// check if file is exist in minio
	objectInfo, err := databases.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return objectInfo.Size > 0
}
