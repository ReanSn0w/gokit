package s3_test

import (
	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/gokit/pkg/client/s3"
)

var (
	client     *s3.Client
	testBucket *s3.BucketMethods
	opts       s3.Options
)

func init() {
	_, err := app.LoadConfiguration("S3 Client test", "debug", &opts)
	if err != nil {
		panic("s3 test init failed: " + err.Error())
	}

	client = opts.MustNewClient()
	testBucket = client.Bucket("gokit-testbucket")
}
