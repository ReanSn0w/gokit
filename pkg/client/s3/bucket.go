package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (c *Client) Bucket(name string) *BucketMethods {
	return &BucketMethods{client: c, name: name}
}

type BucketMethods struct {
	name   string
	client *Client
}

func (bm *BucketMethods) Client() *Client {
	return bm.Client()
}

func (bm *BucketMethods) Get(ctx context.Context, documents ...*GetRequest) GetResponse {
	res := make(GetResponse, len(documents))

	for i := range documents {
		object, err := bm.client.s3.GetObject(ctx, bm.name, documents[i].path, documents[i].opts)

		res[i] = GetResponseItem{
			Path:   documents[i].path,
			Object: object,
			Err:    err,
		}
	}

	return res
}

func (bm *BucketMethods) Put(ctx context.Context, documents ...*PutRequest) PutResponse {
	res := make(PutResponse, len(documents))

	for i := range documents {
		if documents[i] == nil {
			res[i] = PutResponseItem{
				Name:  fmt.Sprintf("item_%v", i),
				Info:  minio.UploadInfo{},
				Error: errors.New("nil put request"),
			}

			continue
		}

		info, err := bm.client.s3.PutObject(
			ctx, bm.name,
			documents[i].name,
			documents[i].reader,
			documents[i].size,
			documents[i].opts)

		res[i] = PutResponseItem{
			Name:  documents[i].name,
			Info:  info,
			Error: err,
		}
	}

	return res
}

func (bm *BucketMethods) Delete(ctx context.Context, documents ...*DeleteRequest) DeleteResponse {
	result := make(DeleteResponse, len(documents))

	for i := range documents {
		result[i] = DeleteResponseItem{
			Path: documents[i].path,
			Err:  bm.client.s3.RemoveObject(ctx, bm.name, documents[i].path, documents[i].opts),
		}
	}

	return result
}
