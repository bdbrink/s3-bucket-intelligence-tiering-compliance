package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"

	"testing"
)

// mock s3 client
type mockS3Client struct {
	s3iface.S3API
}

// response from calling ListBuckets api
func (m *mockS3Client) ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {

	buckets := []*s3.Bucket{
		{
			Name: aws.String("test_bucket1"),
		},
		{
			Name: aws.String("test_bucket2"),
		},
	}
	output := s3.ListBucketsOutput{
		Buckets: buckets,
	}
	return &output, nil
}

// response from calling GetBucketLifecycleConfiguration api
func (m *mockS3Client) GetBucketLifecycleConfiguration(*s3.GetBucketLifecycleConfigurationInput) (*s3.GetBucketLifecycleConfigurationOutput, error) {

	output := &s3.GetBucketLifecycleConfigurationOutput{
		Rules: []*s3.LifecycleRule{
			{
				Filter: &s3.LifecycleRuleFilter{
					Prefix: aws.String(""),
				},
				ID:     aws.String("IntelligentTierRule"),
				Status: aws.String("Enabled"),
				Transitions: []*s3.Transition{
					{
						Days:         aws.Int64(1),
						StorageClass: aws.String("INTELLIGENT_TIERING"),
					},
				},
			},
		},
	}
	return output, nil
}

// unless there is an error this API doesn't return anything https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#PutBucketLifecycleConfigurationOutput
func (m *mockS3Client) PutBucketLifecycleConfiguration(*s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error) {

	output := &s3.PutBucketLifecycleConfigurationOutput{}

	return output, nil
}

// no response from this API https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#PutBucketIntelligentTieringConfigurationOutput
func (m *mockS3Client) PutBucketIntelligentTieringConfiguration(*s3.PutBucketIntelligentTieringConfigurationInput) (*s3.PutBucketIntelligentTieringConfigurationOutput, error) {

	output := &s3.PutBucketIntelligentTieringConfigurationOutput{}

	return output, nil
}

func (m *mockS3Client) GetBucketIntelligentTieringConfiguration(*s3.GetBucketIntelligentTieringConfigurationInput) (*s3.GetBucketIntelligentTieringConfigurationOutput, error) {

	tieringConfig := &s3.IntelligentTieringConfiguration{
		Id:     aws.String("DeepArchiveTest"),
		Status: aws.String("Enabled"),
		Tierings: []*s3.Tiering{
			{
				AccessTier: aws.String("DEEP_ARCHIVE_ACCESS"),
				Days:       aws.Int64(365),
			},
		},
	}
	output := &s3.GetBucketIntelligentTieringConfigurationOutput{
		IntelligentTieringConfiguration: tieringConfig,
	}
	return output, nil
}

func TestListBuckets(t *testing.T) {

	mockClient := &mockS3Client{}
	testListBuckets, _ := listAllBuckets(mockClient)
	t.Run("List buckets should return the 2 buckets from the mock client", func(t *testing.T) {
		assert.Equal(t, testListBuckets, testListBuckets)
	})

}

func TestGetPolicy(t *testing.T) {

	mockClient := &mockS3Client{}
	buckets := []string{"test_bucket1", "test_bucket2"}

	for _, bucket := range buckets {
		testGetPolicy := getPolicy(mockClient, bucket)
		t.Run("buckets should have policy applied to them", func(t *testing.T) {
			assert.Equal(t, testGetPolicy, testGetPolicy)
		})
	}

}

func TestGetTieringPolicy(t *testing.T) {

	mockClient := &mockS3Client{}
	buckets := []string{"test1", "test2"}

	for _, bucket := range buckets {
		testGetTieringPolicy := getTieringPolicy(mockClient, bucket)
		t.Run("apply advanced tiering policy", func(t *testing.T) {
			assert.Equal(t, testGetTieringPolicy, testGetTieringPolicy)
		})
	}

}
