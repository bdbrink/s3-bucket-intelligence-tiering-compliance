package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listBuckets(client *s3.S3) (*s3.ListBucketsOutput, error) {
	res, err := client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getPolicy(client *s3.S3, bucket string) {

	input := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}

	result, err := client.GetBucketLifecycleConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Message() {
			default:
				fmt.Println(aerr.Error())
				fmt.Println(aerr.Message())
			case "The lifecycle configuration does not exist":
				fmt.Printf("policy not found for %v, applying policy \n", bucket)
				putPolicy(client, bucket)
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func putPolicy(client *s3.S3, bucket string) *s3.PutBucketLifecycleConfigurationOutput {

	// after 30 days move objects to intelligent tiering
	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
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
		},
	}

	result, err := client.PutBucketLifecycleConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return result
	}

	fmt.Printf("policy applied to %v \n", bucket)
	return result
}

func main() {

	// default to west for now
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	s3Client := s3.New(sess)

	buckets, err := listBuckets(s3Client)
	if err != nil {
		fmt.Printf("Couldn't list buckets: %v", err)
		return
	}

	for _, bucket := range buckets.Buckets {
		fmt.Printf("Found bucket: %s \n", *bucket.Name)
		getPolicy(s3Client, *bucket.Name)
	}

}
