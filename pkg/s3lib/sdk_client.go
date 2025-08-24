package s3lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type SdkClient struct {
	*s3.Client
}

func NewSdkClient(client *s3.Client) SdkClient {
	return SdkClient{
		Client: client,
	}
}

func (c SdkClient) ConnectionParameters(bucket string) ConnectionParameters {
	var result ConnectionParameters

	ep, err := c.Client.Options().EndpointResolverV2.ResolveEndpoint(context.Background(), s3.EndpointParameters{
		Bucket:         aws.String(bucket),
		Region:         aws.String(c.Client.Options().Region),
		UseFIPS:        aws.Bool(c.Client.Options().EndpointOptions.UseFIPSEndpoint == aws.FIPSEndpointStateEnabled),
		UseDualStack:   aws.Bool(c.Client.Options().EndpointOptions.UseDualStackEndpoint == aws.DualStackEndpointStateEnabled),
		Endpoint:       c.Client.Options().BaseEndpoint,
		ForcePathStyle: aws.Bool(c.Client.Options().UsePathStyle),
		Accelerate:     aws.Bool(c.Client.Options().UseAccelerate),
	})
	if err == nil {
		result.Endpoint = aws.String(ep.URI.String())
	} else {
		result.Endpoint = c.Client.Options().BaseEndpoint
	}

	result.Region = aws.String(c.Client.Options().Region)

	return result
}

func (c SdkClient) ListBuckets(ctx context.Context) Paginator[types.Bucket] {
	bucketPaginator := s3.NewListBucketsPaginator(c.Client, &s3.ListBucketsInput{})
	return ListBucketsPaginator{
		ListBucketsPaginator: bucketPaginator,
	}
}

func (c SdkClient) ListObjects(ctx context.Context, bucket string, prefix string) Paginator[Object] {
	objectPaginator := s3.NewListObjectsV2Paginator(c.Client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	return ListObjectsPaginator{
		ListObjectsV2Paginator: objectPaginator,
	}
}

func (c SdkClient) CreateBucket(ctx context.Context, bucket, region string) error {
	_, err := c.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c SdkClient) UploadFile(ctx context.Context, bucket, key, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = c.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c SdkClient) DownloadFile(ctx context.Context, bucket, key, filePath string) error {
	result, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", key, bucket)
			err = noKey
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucket, key, err)
		}
		return err
	}
	defer result.Body.Close()

	os.Remove(filePath)
	err = os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", filePath, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", key, err)
	}
	_, err = file.Write(body)
	return err
}

func (c SdkClient) GetObject(ctx context.Context, bucket, key string) (ObjectMetadata, error) {
	var result ObjectMetadata
	result.Bucket = bucket
	result.Key = key

	location, err := c.Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		result.Region = string(location.LocationConstraint)
	}

	acl, err := c.Client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		if acl.Owner != nil {
			if acl.Owner.DisplayName != nil {
				result.Owner = acl.Owner.DisplayName
			} else if acl.Owner.ID != nil {
				result.Owner = acl.Owner.ID
			}
		}
	}

	tagging, err := c.Client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		result.Tags = map[string]string{}
		for _, tag := range tagging.TagSet {
			result.Tags[*tag.Key] = *tag.Value
		}
	}

	attr, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		attr.Body.Close()
		result.Type = attr.ContentType
		result.Size = attr.ContentLength
		result.LastModified = attr.LastModified
		result.Metadata = attr.Metadata
		result.ETag = attr.ETag
		result.LegalHold = string(attr.ObjectLockLegalHoldStatus)
	}

	return result, err
}

func ErrorText(err error) string {
	if err == nil {
		return ""
	}

	if smithyErr, ok := findError[*smithy.OperationError](err); ok {
		reason := ""
		if err, ok := findError[*net.OpError](err); ok {
			return err.Error()
		}
		if err, ok := findError[*smithy.GenericAPIError](err); ok {
			reason = err.ErrorMessage()
		}
		if err, ok := findError[*types.BucketAlreadyExists](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.BucketAlreadyOwnedByYou](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.EncryptionTypeMismatch](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.IdempotencyParameterMismatch](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.InvalidObjectState](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.InvalidRequest](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.InvalidWriteOffset](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.NoSuchBucket](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.NoSuchKey](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.NoSuchUpload](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.NotFound](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.ObjectAlreadyInActiveTierError](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.ObjectNotInActiveTierError](err); ok {
			reason = err.Error()
		}
		if err, ok := findError[*types.TooManyParts](err); ok {
			reason = err.Error()
		}

		if reason == "" {
			reason = smithyErr.Error()
		}

		return fmt.Sprintf("%s - %s", smithyErr.Operation(), reason)
	}

	return err.Error()
}

func findError[T any](e error) (T, bool) {
	var errorList []error
	currentErr := e
	for {
		if currentErr == nil {
			break
		}
		errorList = append(errorList, currentErr)
		currentErr = errors.Unwrap(currentErr)
	}

	for _, err := range errorList {
		if e, ok := err.(T); ok {
			return e, true
		}
	}

	var zero T
	return zero, false
}
