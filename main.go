package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

// get all the buckets in the account
func listBuckets(client *s3.S3) (*s3.ListBucketsOutput, error) {
	res, err := client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// if the policy doesn't exist apply it
func getPolicy(client *s3.S3, bucket string) {

	input := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}

	result, err := client.GetBucketLifecycleConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Message() {
			default:
				log.Errorln(aerr.Error())
				log.Errorln(aerr.Message())
			case "The lifecycle configuration does not exist":
				log.Infof("policy not found for %v, applying policy \n", bucket)
				putPolicy(client, bucket)
			}
		} else {
			log.Errorln(err.Error())
		}
	}

	log.Infoln(result)
}

// policy applied if lifecycle not found on the bucket
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
				log.Errorln(aerr.Error())
			}
		} else {
			log.Errorln(err.Error())
		}
	}

	log.Infof("policy applied to %v \n", bucket)
	putTieringPolicy(client, bucket)
	return result
}

// after intelligent tiering applied apply the archive policy
func putTieringPolicy(client *s3.S3, bucket string) *s3.PutBucketIntelligentTieringConfigurationOutput {

	apiObject := &s3.IntelligentTieringConfiguration{
		Id:     aws.String("DeepArchive365"),
		Status: aws.String("Enabled"),
		Tierings: []*s3.Tiering{
			{
				AccessTier: aws.String("DEEP_ARCHIVE_ACCESS"),
				Days:       aws.Int64(365),
			},
		},
	}
	input := &s3.PutBucketIntelligentTieringConfigurationInput{
		Bucket:                          aws.String(bucket),
		Id:                              aws.String("DeepArchive365"),
		IntelligentTieringConfiguration: apiObject,
	}

	result, err := client.PutBucketIntelligentTieringConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Errorln(aerr.Error())
			}
		} else {
			log.Errorln(err.Error())
		}
		return result
	}

	log.Infof("policy applied to %v \n", bucket)
	getTieringPolicy(client, bucket)
	return result
}

// return fresh policy applied to the bucket
func getTieringPolicy(client *s3.S3, bucket string) *s3.GetBucketIntelligentTieringConfigurationOutput {

	input := &s3.GetBucketIntelligentTieringConfigurationInput{
		Bucket: aws.String(bucket),
		Id:     aws.String("DeepArchive365"),
	}

	result, err := client.GetBucketIntelligentTieringConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Message() {
			default:
				log.Errorln(aerr.Error())
				log.Errorln(aerr.Message())
			}
		} else {
			log.Errorln(err.Error())
		}
	}

	log.Infof("Policy found for %v \n", bucket)
	log.Infoln(result.IntelligentTieringConfiguration)
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
		log.Errorf("Couldn't list buckets: %v", err)
		return
	}

	for _, bucket := range buckets.Buckets {
		log.Infof("Found bucket: %s \n", *bucket.Name)
		getPolicy(s3Client, *bucket.Name)
	}

}
