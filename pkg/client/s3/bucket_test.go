package s3_test

import (
	"bytes"
	"context"
	"testing"

	"git.papkovda.ru/library/gokit/pkg/client/s3"
)

var (
	cases = []struct {
		Name   string
		Get    []*s3.GetRequest
		Put    []*s3.PutRequest
		Delete []*s3.DeleteRequest
	}{
		{
			Name: "Работа с одним файлом",
			Get: []*s3.GetRequest{
				s3.NewGetRequest("file.txt"),
			},
			Put: []*s3.PutRequest{
				s3.NewPutRequest(
					bytes.NewReader([]byte("Контент файла 1"))).
					SetName("file.txt"),
			},
			Delete: []*s3.DeleteRequest{
				s3.NewDeleteRequest("file.txt"),
			},
		},
		{
			Name: "Работа с файлом в директории",
			Get: []*s3.GetRequest{
				s3.NewGetRequest("test/file.txt"),
			},
			Put: []*s3.PutRequest{
				s3.NewPutRequest(
					bytes.NewReader([]byte("Контент файла 1"))).
					SetName("test/file.txt"),
			},
			Delete: []*s3.DeleteRequest{
				s3.NewDeleteRequest("test/file.txt"),
			},
		},
		{
			Name: "Работа с файлом без имени",
			Put: []*s3.PutRequest{
				s3.NewPutRequest(
					bytes.NewReader([]byte("Контент файла 1"))),
			},
		},
	}
)

func Test_BucketMethods(t *testing.T) {
	t.Run("put", testBucketMethods_Put)
	t.Run("get", testBuckerMethods_Get)
	t.Run("delete", testBucketMethods_Delete)
}

func testBuckerMethods_Get(t *testing.T) {
	for _, c := range cases {
		t.Log("[INFO] запуск кейса:", c.Name)
		response := testBucket.Get(context.Background(), c.Get...)
		if err := response.Err(); err != nil {
			t.Logf("[ERROR] при запуске кейса %v произошли следующие ошибки: %v", c.Name, err)
			t.Error(err)
		}

		if err := response.Close(); err != nil {
			t.Logf("[ERROR] при закрытии файлов в кейсе %v произошла ошибка: %v", c.Name, err)
			t.Error(err)
		}
	}
}

func testBucketMethods_Put(t *testing.T) {
	for _, c := range cases {
		t.Log("[INFO] запуск кейса:", c.Name)
		response := testBucket.Put(context.Background(), c.Put...)
		if err := response.Err(); err != nil {
			t.Logf("[ERROR] при запуске кейса %v произошли следующие ошибки: %v", c.Name, err)
			t.Error(err)
		}
	}
}

func testBucketMethods_Delete(t *testing.T) {
	for _, c := range cases {
		t.Log("[INFO] запуск кейса:", c.Name)
		response := testBucket.Delete(context.Background(), c.Delete...)
		if err := response.Err(); err != nil {
			t.Logf("[ERROR] при запуске кейса %v произошли следующие ошибки: %v", c.Name, err)
			t.Error(err)
		}
	}
}
