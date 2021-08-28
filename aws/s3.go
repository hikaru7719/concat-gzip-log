package aws

import (
	"bytes"
	"context"
	"io"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type SortableObjects []types.Object

func (s SortableObjects) Len() int {
	return len(s)
}

func (s SortableObjects) Less(i, j int) bool {
	return s[i].LastModified.Before(*s[j].LastModified)
}

func (s SortableObjects) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StorageReader struct {
	client Client
	bucket string
	date   time.Time
}

func NewStorageReader(bucket string, date time.Time) *StorageReader {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {

	}

	return &StorageReader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		date:   date,
	}
}

func (s *StorageReader) IsValidDateRange(lastModify *time.Time) bool {
	y, m, d := s.date.Date()
	ty, tm, td := lastModify.Date()
	return y == ty && m == tm && d == td
}

func (s *StorageReader) FilterByDate(objects []types.Object) []types.Object {
	result := make([]types.Object, 0)
	for _, o := range objects {
		if s.IsValidDateRange(o.LastModified) {
			result = append(result, o)
		}
	}
	return result
}

func (s *StorageReader) SortByLastModify(objects []types.Object) []types.Object {
	sort.Sort(SortableObjects(objects))
	return objects
}

func (s *StorageReader) MapKey(objects []types.Object) []string {
	result := make([]string, 0)
	for _, o := range objects {
		result = append(result, *o.Key)
	}
	return result
}

func (s *StorageReader) Filter(objects []types.Object) []string {
	return s.MapKey(s.SortByLastModify(s.FilterByDate(objects)))
}

func (s *StorageReader) List() []string {
	out, err := s.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		log.Print(err)
		return []string{}
	}
	return s.Filter(out.Contents)
}

func (s *StorageReader) GetObject(key string) bytes.Buffer {
	out, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Print(err)
	}
	return s.CopyBuffer(out.Body)
}

func (s *StorageReader) CopyBuffer(body io.ReadCloser) bytes.Buffer {
	var buf bytes.Buffer
	defer body.Close()
	io.Copy(&buf, body)
	return buf
}

func (s *StorageReader) GetAllObject(keys []string) []bytes.Buffer {
	result := make([]bytes.Buffer, 0)
	for _, k := range keys {
		result = append(result, s.GetObject(k))
	}
	return result
}
