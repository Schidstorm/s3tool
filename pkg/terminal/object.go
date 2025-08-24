package terminal

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

func editObject(c Context) error {
	tmpFilePath, err := downloadFileToTmp(c)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	oldHash, err := fileHash(tmpFilePath)
	if err != nil {
		return err
	}

	err = EditFile(c, tmpFilePath)
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
		err = c.S3Client().UploadFile(context.Background(), c.Bucket(), c.ObjectKey(), tmpFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func viewObject(c Context) error {
	tmpFilePath, err := downloadFileToTmp(c)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	return ShowFile(c, tmpFilePath)
}

func downloadFileToTmp(c Context) (string, error) {
	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		return "", err
	}
	tmpFilePath := tmpDir + "/" + c.ObjectKey()

	err = c.S3Client().DownloadFile(context.Background(), c.Bucket(), c.ObjectKey(), tmpFilePath)
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
