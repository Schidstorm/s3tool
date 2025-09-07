package terminal

import (
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/internal/s3lib"
)

type ModalBuilder func(close func()) tview.Primitive

type Context interface {
	S3Client() s3lib.Client
	Bucket() string
	ObjectKey() string
	Modal(build ModalBuilder)
	SetError(err error)
	OpenPage(page PageContent)
	SuspendApp(f func()) bool

	WithClient(client s3lib.Client) Context
	WithBucket(bucket string) Context
	WithObjectKey(key string) Context
	WithModalFunc(f func(build ModalBuilder)) Context
	WithErrorFunc(f func(err error)) Context
	WithOpenPageFunc(f func(page PageContent)) Context
	WithSuspendAppFunc(f func(func()) bool) Context
}

type contextImpl struct {
	client     s3lib.Client
	bucket     string
	objectKey  string
	modalFunc  func(build ModalBuilder)
	errorFunc  func(err error)
	openFunc   func(page PageContent)
	suspendApp func(func()) bool
}

func NewContext() Context {
	return contextImpl{}
}

func (c contextImpl) S3Client() s3lib.Client {
	return c.client
}

func (c contextImpl) Bucket() string {
	return c.bucket
}

func (c contextImpl) ObjectKey() string {
	return c.objectKey
}

func (c contextImpl) Modal(build ModalBuilder) {
	if c.modalFunc != nil {
		c.modalFunc(build)
	}
}

func (c contextImpl) SetError(err error) {
	if c.errorFunc != nil {
		c.errorFunc(err)
	}
}

func (c contextImpl) OpenPage(page PageContent) {
	if c.openFunc != nil {
		c.openFunc(page)
	}
}

func (c contextImpl) SuspendApp(f func()) bool {
	if c.suspendApp != nil {
		return c.suspendApp(f)
	}
	return false
}

func (c contextImpl) WithClient(client s3lib.Client) Context {
	c.client = client
	return c
}

func (c contextImpl) WithBucket(bucket string) Context {
	c.bucket = bucket
	return c
}

func (c contextImpl) WithObjectKey(key string) Context {
	c.objectKey = key
	return c
}

func (c contextImpl) WithModalFunc(f func(build ModalBuilder)) Context {
	c.modalFunc = f
	return c
}

func (c contextImpl) WithErrorFunc(f func(err error)) Context {
	c.errorFunc = f
	return c
}

func (c contextImpl) WithOpenPageFunc(f func(page PageContent)) Context {
	c.openFunc = f
	return c
}

func (c contextImpl) WithSuspendAppFunc(f func(func()) bool) Context {
	c.suspendApp = f
	return c
}
