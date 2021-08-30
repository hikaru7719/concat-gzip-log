package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func testParseDate(t *testing.T, dateString string) *time.Time {
	t.Helper()
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		t.Fatal(err)
	}
	return &date
}

type MockClient struct{}

func (m *MockClient) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key:          aws.String("valid2"),
				LastModified: aws.Time(time.Date(2021, time.August, 28, 10, 1, 0, 0, time.UTC)),
			},
			{
				Key:          aws.String("valid"),
				LastModified: aws.Time(time.Date(2021, time.August, 28, 10, 0, 0, 0, time.UTC)),
			},
			{
				Key:          aws.String("invalid"),
				LastModified: aws.Time(time.Date(2021, time.August, 29, 10, 0, 0, 0, time.UTC)),
			},
			{
				Key:          aws.String("valid3"),
				LastModified: aws.Time(time.Date(2021, time.August, 28, 10, 2, 0, 0, time.UTC)),
			},
		},
	}, nil
}

func (m *MockClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return nil, nil
}

func newMockStorageReader(bucket string, date time.Time) *StorageReader {
	return &StorageReader{
		client: &MockClient{},
		bucket: bucket,
		date:   date,
	}
}

func TestIsValidDateRange(t *testing.T) {
	cases := map[string]struct {
		validDate time.Time
		target    time.Time
		expect    bool
	}{
		"true": {
			validDate: *testParseDate(t, "2021-08-28"),
			target:    *testParseDate(t, "2021-08-28"),
			expect:    true,
		},
		"false": {
			validDate: *testParseDate(t, "2021-08-28"),
			target:    *testParseDate(t, "2021-08-29"),
			expect:    false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			client := newMockStorageReader("", tc.validDate)
			actual := client.IsValidDateRange(&tc.target)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestList(t *testing.T) {
	cases := map[string]struct {
		date   time.Time
		expect []string
	}{
		"ListObject": {
			date:   *testParseDate(t, "2021-08-28"),
			expect: []string{"valid", "valid2", "valid3"},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			client := newMockStorageReader("", tc.date)
			actual := client.List()
			assert.Equal(t, tc.expect, actual)
		})
	}
}
