package s3

import (
	"io"

	"github.com/ReanSn0w/gokit/pkg/tool"
	"github.com/minio/minio-go/v7"
)

func NewGetRequest(path string) *GetRequest {
	return &GetRequest{path: path}
}

type GetRequest struct {
	path string
	opts minio.GetObjectOptions
}

func (g *GetRequest) SetOptions(configure func(*minio.GetObjectOptions)) *GetRequest {
	configure(&g.opts)
	return g
}

type GetResponse []GetResponseItem

func (g GetResponse) Err() error {
	es := tool.NewErrorsMap()

	for _, f := range g {
		if err := f.Err; err != nil {
			es[f.Path] = err
		}
	}

	return es.IsError()
}

func (g GetResponse) Close() error {
	es := tool.NewErrorsMap()

	for _, f := range g {
		if f.Err != nil {
			continue
		}

		err := f.Object.Close()
		if err != nil {
			es[f.Path] = err
		}
	}

	return es.IsError()
}

type GetResponseItem struct {
	Path   string
	Err    error
	Object *minio.Object
}

func NewPutRequest(reader io.Reader) *PutRequest {
	return &PutRequest{
		name:   tool.NewID() + ".bin",
		reader: reader,
		size:   -1,
		opts:   minio.PutObjectOptions{},
	}
}

type PutRequest struct {
	name   string
	reader io.Reader
	size   int64
	opts   minio.PutObjectOptions
}

func (p *PutRequest) SetName(name string) *PutRequest {
	p.name = name
	return p
}

func (p *PutRequest) SetSize(size int64) *PutRequest {
	p.size = size
	return p
}

func (p *PutRequest) SetOptions(configure func(*minio.PutObjectOptions)) *PutRequest {
	configure(&p.opts)
	return p
}

type PutResponse []PutResponseItem

func (p PutResponse) Err() error {
	es := tool.NewErrorsMap()

	for i := range p {
		if p[i].Error != nil {
			es[p[i].Name] = p[i].Error
		}
	}

	return es.IsError()
}

type PutResponseItem struct {
	Name  string
	Info  minio.UploadInfo
	Error error
}

func NewDeleteRequest(path string) *DeleteRequest {
	return &DeleteRequest{path: path}
}

type DeleteRequest struct {
	path string
	opts minio.RemoveObjectOptions
}

type DeleteResponse []DeleteResponseItem

func (d DeleteResponse) Err() error {
	es := tool.NewErrorsMap()

	for _, item := range d {
		if err := item.Err; err != nil {
			es[item.Path] = err
		}
	}

	return es.IsError()
}

type DeleteResponseItem struct {
	Path string
	Err  error
}
