package s3lib

import "github.com/aws/aws-sdk-go-v2/service/s3"

var active3Client *s3.Client

type ActiveClientChangeHander func(*s3.Client)

var activeClientChangeHandlerRegistry []ActiveClientChangeHander

func GetActiveClient() (*s3.Client, bool) {
	if active3Client == nil {
		return nil, false
	}
	return active3Client, true
}

func SetActiveClient(client *s3.Client) {
	active3Client = client
	for _, handler := range activeClientChangeHandlerRegistry {
		handler(client)
	}
}

func RegisterActiveClientChangeHandler(handler ActiveClientChangeHander) {
	activeClientChangeHandlerRegistry = append(activeClientChangeHandlerRegistry, handler)
}
