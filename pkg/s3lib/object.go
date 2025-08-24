package s3lib

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Object struct {
	Kind   ObjectKind
	Object types.Object
}

func NewObjectFile(obj types.Object) Object {
	return Object{
		Kind:   ObjectKindFile,
		Object: obj,
	}
}

func NewObjectDirectory(name string) Object {
	return Object{
		Kind: ObjectKindDirectory,
		Object: types.Object{
			Key: aws.String(name),
		},
	}
}

type ObjectKind int

const (
	ObjectKindFile ObjectKind = iota
	ObjectKindDirectory
)

func (o Object) IsFile() bool {
	return o.Kind == ObjectKindFile
}
func (o Object) IsDirectory() bool {
	return o.Kind == ObjectKindDirectory
}
