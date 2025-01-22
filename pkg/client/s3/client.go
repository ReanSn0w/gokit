package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(tls bool, endpoint, accessKey, secretKey string) (*Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: tls,
	})

	return &Client{s3: client}, err
}

type Client struct {
	s3 *minio.Client
}

type Options struct {
	S3 struct {
		Endpoint  string `long:"endpoint" env:"ENDPOINT" description:"s3 storage endpoint"`
		AccessKey string `long:"access-key" env:"ACCESS_KEY" description:"s3 access key"`
		SecretKey string `long:"secret-key" env:"SECRET_KEY" description:"s3 secret key"`
		TLS       bool   `long:"tls" env:"TLS" description:"s3 tls connection"`
	} `group:"s3" namespace:"s3" env-namespace:"S3"`
}

func (o *Options) NewClient() (*Client, error) {
	return New(o.S3.TLS, o.S3.Endpoint, o.S3.AccessKey, o.S3.SecretKey)
}

func (o *Options) MustNewClient() *Client {
	client, err := o.NewClient()
	if err != nil {
		panic("s3 client setup failed: " + err.Error())
	}

	return client
}
