package terminal

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"os"

	"github.com/schidstorm/s3tool/pkg/s3lib"
)

func editObject(client s3lib.Client, bucketName, objectName string) error {
	tmpFilePath, err := downloadFileToTmp(client, bucketName, objectName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	oldHash, err := fileHash(tmpFilePath)
	if err != nil {
		return err
	}

	err = EditFile(tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		return errors.New("file does not exist after editing")
	}

	newHash, err := fileHash(tmpFilePath)
	if err != nil {
		return err
	}

	if oldHash != newHash {
		err = client.UploadFile(context.Background(), bucketName, objectName, tmpFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func viewObject(client s3lib.Client, bucketName, objectName string) error {
	tmpFilePath, err := downloadFileToTmp(client, bucketName, objectName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	return ShowFile(tmpFilePath)
}

func downloadFileToTmp(client s3lib.Client, bucketName, objectName string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		return "", err
	}
	tmpFilePath := tmpDir + "/" + objectName

	err = client.DownloadFile(context.Background(), bucketName, objectName, tmpFilePath)
	if err != nil {
		return "", err
	}

	return tmpFilePath, nil

}

func fileHash(filePath string) (string, error) {
	sha256Hash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buffer [4096]byte
	for {
		n, err := file.Read(buffer[:])
		if n > 0 {
			sha256Hash.Write(buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	return string(sha256Hash.Sum(nil)), nil
}
